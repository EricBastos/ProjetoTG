package operator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/configs"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractCaller"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractInterface"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/rabbitmqClient"
	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Operator struct {
	rabbit              *rabbitmqClient.RabbitMQClient
	smartContractOpRepo database.SmartcontractOperationInterface
	feedbackRepo        database.FeedbackInterface
	transferDb          database.TransferInterface
	burnOpDb            database.BurnOpInterface
	transactionNotify   chan *contractInterface.OnGoingTransactionData
	interrupt           chan os.Signal
	client              *ethclient.Client
	online              bool
	notifyOffline       chan bool
	mu                  sync.RWMutex
}

type PublishData struct {
	respJson []byte
	op       string
}

type TxWithError struct {
	SmartcontractCallId string
	OperationName       string
	Tx                  *contractInterface.TxResult
	Error               error
}

type MintData struct {
	Id                          string `json:"id"`
	Chain                       string `json:"chain"`
	WalletAddress               string `json:"address"`
	Amount                      int    `json:"amount"`
	AssociatedBankTransactionId string `json:"associatedBankTransactionId"`
}

type BurnData struct {
	Id            string `json:"id"`
	Chain         string `json:"chain"`
	WalletAddress string `json:"address"`
	Amount        int    `json:"amount"`
	UserName      string `json:"userName"`
	UserTaxId     string `json:"userTaxId"`
	AccBankCode   string `json:"accBankCode"`
	AccBranchCode string `json:"accBranchCode"`
	AccNumber     string `json:"accNumber"`
	//AssociatedBankTransactionId string `json:"associatedBankTransactionId"`
	Permit *entities.PermitData `json:"permit"`
}

type FeedbackResponse struct {
	ID                  string `json:"id"`
	OperationId         string `json:"operationId"`
	OperationOriginType string `json:"operationOriginType"`
	Success             bool   `json:"success"`
	ErrorMsg            string `json:"errorMsg"`
}

type GenericMessage struct {
	ID                  string      `json:"id"`
	IsRetry             bool        `json:"isRetry"`
	Operation           string      `json:"operation"`
	ResponsibleUser     string      `json:"userId"`
	OperationOriginType string      `json:"operationOriginType"`
	WorkspaceId         string      `json:"workspaceId"`
	Data                interface{} `json:"data"`
}

func NewOperator(smartcontractOpDb database.SmartcontractOperationInterface, feedbackRepo database.FeedbackInterface, burnOpDb database.BurnOpInterface, transferDb database.TransferInterface, rabbit *rabbitmqClient.RabbitMQClient, notifyOff chan bool) *Operator {
	op := &Operator{}
	op.feedbackRepo = feedbackRepo
	op.burnOpDb = burnOpDb
	op.transferDb = transferDb
	op.notifyOffline = notifyOff
	op.setOnlineStatus(true)
	op.smartContractOpRepo = smartcontractOpDb
	op.transactionNotify = make(chan *contractInterface.OnGoingTransactionData, 1000)
	op.rabbit = rabbit
	return op
}

func (o *Operator) Start() {

	defer o.setOnlineStatus(false)

	ethConfig := contractInterface.EthContractConfig{
		RPCHost:           configs.Cfg.EthereumRpcHost,
		ContractAddress:   configs.Cfg.EthereumTokenContract,
		OwnerPrivateKey:   configs.Cfg.EthereumWalletPrivateKey,
		TransactionNotify: o.transactionNotify,
	}

	polyConfig := contractInterface.PolygonContractConfig{
		RPCHost:           configs.Cfg.PolygonRpcHost,
		ContractAddress:   configs.Cfg.PolygonTokenContract,
		OwnerPrivateKey:   configs.Cfg.PolygonWalletPrivateKey,
		TransactionNotify: o.transactionNotify,
	}

	contractCaller.InitializeContractCaller(
		&ethConfig,
		&polyConfig,
	)

	go o.trackTransactions()

	consumeWg := sync.WaitGroup{}

	for chain := range o.rabbit.Messages {
		consumeWg.Add(1)
		go o.consumeChanLoop(chain, &consumeWg)
	}

	consumeWg.Wait()
}

func (o *Operator) consumeChanLoop(chain string, wg *sync.WaitGroup) {
	defer wg.Done()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
MainForLoop:
	for {
		msgChan := o.rabbit.GetMsgChan(chain)
		if msgChan == nil {
			continue
		}
		select {
		case m, ok := <-msgChan:
			startTime := time.Now().UnixMilli()
			if !ok {
				o.rabbit.SetMsgChan(chain, nil)
				log.Println("Rabbit channel closed. On hold until reopen", m)
				continue
			}
			if m.Redelivered {
				log.Println("Message redelivered. It most likely it's already been processed, so ignored.", string(m.Body))
				_ = o.rabbit.Ch.Ack(m.DeliveryTag, false)
				continue
			}
			{
				var message GenericMessage
				var err error
				var txsWithError []TxWithError
				err = json.Unmarshal(m.Body, &message)
				if err != nil {
					log.Println("Error unmarshalling message. Command ignored")
					o.nackMsg("", m.DeliveryTag, "", "", false, "", "Error unmarshalling message: "+string(m.Body), false)
					continue
				}
				log.Println("Received message:")
				log.Println(message)
				dataJson, err := json.Marshal(message.Data)
				log.Println(string(dataJson))
				if err != nil {
					log.Println("Error encoding operation data", err)
					o.nackMsg(message.OperationOriginType, m.DeliveryTag, message.Operation, message.ID, false, "", "Error encoding operation data: "+string(m.Body), message.IsRetry)
					continue
				}
				ongoingTransaction := &contractInterface.OnGoingTransactionData{
					Operation:           message.Operation,
					ID:                  message.ID,
					OperationOriginType: message.OperationOriginType,
				}
				log.Println(ongoingTransaction)
				switch message.Operation {
				case "MINT":
					var data MintData
					err = json.Unmarshal(dataJson, &data)
					if err != nil {
						log.Println("Error unmarshalling MINT data. Command ignored", err)
						o.nackMsg(message.OperationOriginType, m.DeliveryTag, message.Operation, message.ID, false, "", "Error unmarshalling MINT data", message.IsRetry)
						continue
					}
					txsWithError = o.Mint(&data, ongoingTransaction)
				case "BURN":
					var data BurnData
					err = json.Unmarshal(dataJson, &data)
					if err != nil {
						log.Println("Error unmarshalling BURN data. Command ignored", err)
						o.nackMsg(message.OperationOriginType, m.DeliveryTag, message.Operation, message.ID, false, "", "Error unmarshalling BURN data", message.IsRetry)
						continue
					}
					txsWithError = o.Burn(&data, ongoingTransaction)
				}

				err2 := o.rabbit.Ch.Ack(m.DeliveryTag, false)
				if err2 != nil {
					log.Println("Error: could not ack message.")
				}

				for _, res := range txsWithError {
					tx := res.Tx
					err := res.Error

					if err != nil {
						log.Println("Operation failed:", message.Operation, err.Error())
						operation := entities.NewSmartcontractOperation(message.Operation, message.OperationOriginType, message.ID, false, "", err.Error(), message.IsRetry)
						o.writeOpInDB(operation)

						resp := FeedbackResponse{
							ID:                  operation.ID.String(),
							OperationId:         message.ID,
							OperationOriginType: message.OperationOriginType,
							Success:             false,
							ErrorMsg:            err.Error(),
						}

						o.publishResult(resp, message.Operation)
					} else {
						log.Println("Posted transaction:", tx.Hash)
						operation := entities.NewSmartcontractOperation(res.OperationName, message.OperationOriginType, message.ID, true, tx.Hash, "", message.IsRetry)
						parsedId, _ := entities2.ParseID(res.SmartcontractCallId)
						operation.ID = &parsedId
						o.writeOpInDB(operation)
					}

				}
			}
			endTime := time.Now().UnixMilli()
			remainingInterval := 200 - (endTime - startTime)
			if remainingInterval > 0 {
				time.Sleep(time.Duration(remainingInterval) * time.Millisecond)
			}

		case <-interrupt:
			log.Println("Ctrl-C Detected. Shuting down chain", chain)
			break MainForLoop
		}
	}
}

func errorsToRepeat(err error) bool {
	if strings.Contains(err.Error(), "429") ||
		strings.Contains(err.Error(), "nonce too low") ||
		strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "cannot estimate") ||
		strings.Contains(err.Error(), "transaction underpriced") {
		return true
	}
	return false
}

func (o *Operator) Mint(data *MintData, txData *contractInterface.OnGoingTransactionData) []TxWithError {
	txData.WalletAddress = data.WalletAddress

	var resp *contractInterface.TxResult
	var respId string
	operation := func() error {
		txData.PostedTime = time.Now().UnixMilli()
		txData.SmartcontractCallId = entities2.NewID().String()
		res, err := contractCaller.CC.Mint(txData, false, data.Chain, data.WalletAddress, data.Amount)
		if err != nil {
			log.Println("(MINT) Backoff try failed:", err.Error())
			if errorsToRepeat(err) {
				return err
			}
			return backoff.Permanent(err)
		}
		resp = res
		respId = txData.SmartcontractCallId
		return nil
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     1 * time.Second,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	err := backoff.Retry(operation, b)
	return []TxWithError{
		{
			OperationName:       txData.Operation,
			SmartcontractCallId: respId,
			Tx:                  resp,
			Error:               err,
		},
	}
}

func (o *Operator) Burn(data *BurnData, txData *contractInterface.OnGoingTransactionData) []TxWithError {
	txData.WalletAddress = data.WalletAddress

	var resp *contractInterface.TxResult
	var respId string
	operation := func() error {
		txData.PostedTime = time.Now().UnixMilli()
		txData.SmartcontractCallId = entities2.NewID().String()
		res, err := contractCaller.CC.Burn(txData, false, data.Chain, data.WalletAddress, data.Amount, data.Permit)
		if err != nil {
			log.Println("(BURN) Backoff try failed:", err.Error())
			if errorsToRepeat(err) {
				return err
			}
			return backoff.Permanent(err)
		}
		resp = res
		respId = txData.SmartcontractCallId
		return nil
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     1 * time.Second,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	err := backoff.Retry(operation, b)
	//if err != nil {
	//	_ = globals.Cache.Delete(context.Background(), "permit:"+data.WalletAddress)
	//}
	return []TxWithError{
		{
			OperationName:       txData.Operation,
			SmartcontractCallId: respId,
			Tx:                  resp,
			Error:               err,
		},
	}
}

func (o *Operator) trackTransactions() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case info := <-o.transactionNotify:
			log.Printf("%+v\n", info)

			resp := FeedbackResponse{
				ID:                  info.SmartcontractCallId,
				OperationId:         info.ID,
				OperationOriginType: info.OperationOriginType,
				Success:             info.Success,
				ErrorMsg:            info.ErrReason,
			}

			o.publishResult(resp, info.Operation)
		case <-interrupt:
			return
		}
	}
}

func (o *Operator) setOnlineStatus(val bool) {
	o.mu.Lock()
	o.online = val
	o.mu.Unlock()
	if val == false {
		o.notifyOffline <- true
	}
}

func (o *Operator) GetOnlineStatus() bool {
	var result bool
	o.mu.RLock()
	result = o.online
	o.mu.RUnlock()
	return result
}

func (o *Operator) publishResult(resp FeedbackResponse, op string) {
	feedback := entities.NewFeedback(
		resp.OperationId,
		op,
		resp.ID,
		resp.Success,
		resp.ErrorMsg,
	)
	err := o.feedbackRepo.Create(feedback)
	if err != nil {
		feedbackData := fmt.Sprintf("%+v", feedback)
		log.Println("(Blockchain Feedback) Couldn't persist blockchain feedback. Data: " + feedbackData + ", err: " + err.Error())
	}

	// If it's a successful burn, we follow up by creating a mocked bank transfer
	if op == "BURN" && resp.Success {
		go func() {
			err := o.mockBankTransfer(resp)
			if err != nil {
				log.Println("(SANDBOX) Error mocking bank transfer: " + err.Error())
			}
		}()
	}
}

func (o *Operator) nackMsg(origin string, tag uint64, op, id string, exec bool, tx, reason string, isRetry bool) {
	err := o.rabbit.Ch.Nack(tag, false, false)
	if err != nil {
		log.Println("Error: could not nack message.")
	} else {
		operation := entities.NewSmartcontractOperation(op, origin, id, exec, tx, reason, isRetry)
		go o.writeOpInDB(operation)
	}
}

func (o *Operator) writeOpInDB(operation *entities.SmartcontractOperation) {
	//log.Println("Sending to DB:", op, id, data, executed, tx, reason)
	err := o.smartContractOpRepo.Create(operation)
	if err != nil {
		log.Println("Error writing op to DB:", err.Error())
	}
	//log.Println("Sent to DB")
}

func (o *Operator) mockBankTransfer(feedback FeedbackResponse) error {

	burnOp, err := o.burnOpDb.Get(feedback.OperationId)
	if err != nil {
		return err
	}

	transf := entities.NewTransfer(
		burnOp.WalletAddress,
		burnOp.Amount,
		burnOp.UserName,
		burnOp.UserTaxId,
		burnOp.AccBankCode,
		burnOp.AccBranchCode,
		burnOp.AccNumber,
		burnOp.ResponsibleUser,
		burnOp.Chain,
		burnOp.Id.String(),
	)
	err = o.transferDb.Create(transf)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	body, _ := json.Marshal(map[string]interface{}{
		"subscription": "transfer",
		"data": map[string]interface{}{
			"transferId": transf.Id,
		},
	})

	webhookUrl := "http://" + configs.Cfg.BankWebhookHost + ":" + configs.Cfg.BankWebhookPort
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewReader(body))
	req.Close = true
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
