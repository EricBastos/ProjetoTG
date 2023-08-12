package handlers

import (
	"encoding/json"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/dtos"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/usecases"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"io"
	"log"
	"net/http"
)

type EventWrapper struct {
	Subscription string          `json:"subscription"`
	Data         json.RawMessage `json:"data"`
}

type WebhookHandler struct {
	mintOperationDb         database.MintOpInterface
	transferDb              database.TransferInterface
	transferFeedbackDb      database.TransferFeedbackInterface
	staticDepositDb         database.StaticDepositInterface
	staticDepositFeedbackDb database.StaticDepositFeedbackInterface
	rabbitClient            *rabbitmqClient.RabbitMQClient
}

func NewWebhookHandler(
	mintOperationDb database.MintOpInterface,
	transferDb database.TransferInterface,
	transferFeedbackDb database.TransferFeedbackInterface,
	staticDepositDb database.StaticDepositInterface,
	staticDepositFeedbackDb database.StaticDepositFeedbackInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,

) *WebhookHandler {
	return &WebhookHandler{
		mintOperationDb:         mintOperationDb,
		transferDb:              transferDb,
		transferFeedbackDb:      transferFeedbackDb,
		staticDepositDb:         staticDepositDb,
		staticDepositFeedbackDb: staticDepositFeedbackDb,
		rabbitClient:            rabbitClient,
	}
}

func (h *WebhookHandler) Listen(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	var event EventWrapper

	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Println("(Webhook) Couldn't unmarshal received event data: " + err.Error() + ", data: " + string(data))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Println("Received", event.Subscription, "subscription")
	switch event.Subscription {
	//case "transfer":
	//	var transfer dtos.TransferFeedbackInput
	//	err = json.Unmarshal(logData, &transfer)
	//	if err != nil {
	//		logger.Critical("(Webhook) Couldn't unmarshal received transfer data: " + err.Error() + ", data: " + string(logData))
	//		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//	if _, ok := utils.IgnoreList[utils.TrimCpfCnpj(transfer.Transfer.TaxId)]; !ok {
	//		usecase := usecases.NewProcessTransferFeedbackUsecase(
	//			h.mintOperationDb,
	//			h.transferDb,
	//			h.rabbit,
	//			h.transferFeedbackDb,
	//			h.pendingTransferDb,
	//			h.operationsToRetryDb,
	//			h.emailNotifier,
	//			h.discordBot,
	//		)
	//		err = usecase.Process(&transfer, event.Event.WorkspaceId)
	//		if err != nil {
	//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//			return
	//		}
	//	}

	case "deposit":
		var deposit dtos.DepositFeedbackInput
		err = json.Unmarshal(event.Data, &deposit)
		if err != nil {
			log.Println("(Webhook) Couldn't unmarshal received deposit data: ", err.Error(), " data: ", deposit)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		usecase := usecases.NewProcessDepositFeedbackUsecase(
			h.mintOperationDb,
			h.staticDepositDb,
			h.staticDepositFeedbackDb,
			h.rabbitClient,
		)
		err = usecase.Process(&deposit)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
