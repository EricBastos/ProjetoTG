package rabbitmqClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

type RabbitMQClient struct {
	mu       sync.Mutex
	Messages map[string]<-chan amqp.Delivery
	confirms chan amqp.Confirmation
	Ch       *amqp.Channel
	conn     *amqp.Connection
	config   *struct {
		user             string
		pass             string
		host             string
		port             string
		consumerQueues   map[string]string
		consumerExchange string
		producerExchange string
	}
}

func NewRabbitMQClient(user, pass, host, port, consumerExchange, producerExchange string, consumerQueues map[string]string) (*RabbitMQClient, error) {
	l := &RabbitMQClient{config: &struct {
		user             string
		pass             string
		host             string
		port             string
		consumerQueues   map[string]string
		consumerExchange string
		producerExchange string
	}{
		user:             user,
		pass:             pass,
		host:             host,
		port:             port,
		consumerQueues:   consumerQueues,
		consumerExchange: consumerExchange,
		producerExchange: producerExchange,
	},
		Messages: map[string]<-chan amqp.Delivery{}}
	err := l.initialize()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *RabbitMQClient) initialize() error {
	rabbitmqConnString := fmt.Sprintf("amqp://%s:%s@%s:%s/vhost",
		l.config.user,
		l.config.pass,
		l.config.host,
		l.config.port,
	)

	var err error
	l.conn, err = amqp.Dial(rabbitmqConnString)
	if err != nil {
		log.Println("Dial failed")
		return err
	}
	l.Ch, err = l.conn.Channel()
	if err != nil {
		return err
	}

	err = l.Ch.Confirm(false)
	if err != nil {
		return err
	}

	l.confirms = l.Ch.NotifyPublish(make(chan amqp.Confirmation, 10000)) // Number of expected concurrent requests

	err = l.setupConsumer()
	if err != nil {
		return err
	}

	err = l.setupProducer()
	if err != nil {
		return err
	}

	go l.observeConnection()
	return nil
}

func (l *RabbitMQClient) observeConnection() {
	connErr := <-l.conn.NotifyClose(make(chan *amqp.Error, 2))
	log.Printf("Lost rabbitmq connection due to error: %s\n", connErr)
	for {
		l.closeConnections()
		if err := l.initialize(); err == nil {
			log.Println("Connection to rabbitmq reestablished")
			go l.observeConnection()
			return
		}
		//log.Println("Failed to reconnect to rabbitmq. Retrying in 5 seconds")
		time.Sleep(30 * time.Second)
	}
}

func (l *RabbitMQClient) closeConnections() {
	if l.Ch != nil {
		_ = l.Ch.Close()
	}
	if l.conn != nil && l.conn.IsClosed() {
		_ = l.conn.Close()
	}
}

func (l *RabbitMQClient) setupConsumer() error {

	err := l.Ch.ExchangeDeclare(
		l.config.consumerExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for c, q := range l.config.consumerQueues {
		args := amqp.Table{ // queue args
			"x-max-priority": 3,
		}
		commandQueue, err := l.Ch.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			args,  // args
		)
		if err != nil {
			return err
		}

		err = l.Ch.QueueBind(q, c, l.config.consumerExchange, false, nil)
		if err != nil {
			return err
		}
		err = l.Ch.Qos(1, 0, false)
		if err != nil {
			return err
		}
		msgQ, err := l.Ch.Consume(
			commandQueue.Name, // queue
			"",                // consumer
			false,             // auto-ack
			false,             // exclusive
			false,             // no-local
			false,             // no-wait
			nil,
		)
		if err != nil {
			return err
		}
		l.SetMsgChan(c, msgQ)
	}
	//log.Println("Msg chan:", l.Messages)
	return nil
}

func (l *RabbitMQClient) setupProducer() error {
	err := l.Ch.ExchangeDeclare(
		l.config.producerExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (l *RabbitMQClient) SecurePublish(data []byte, routingKey string, headers map[string]interface{}, exchange string, priority uint8) error {
	conf, err := l.Ch.PublishWithDeferredConfirmWithContext(
		context.Background(),
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			Priority:     priority,
		})
	if err != nil {
		return err
	}

	confirmation := conf.Wait()

	if confirmation {
		return nil
	} else {
		return errors.New("rabbitmq didn't ack publishing")
	}
}

func (l *RabbitMQClient) CallSmartcontract(op entities.SmartContractOp, opOriginType entities.OperationOriginType, priority uint8) error {

	type OpToSendStruct struct {
		ID                  string      `json:"id"`
		IsRetry             bool        `json:"isRetry"`
		UserId              string      `json:"userId"`
		WorkspaceId         string      `json:"workspaceId"`
		OperationOriginType string      `json:"operationOriginType"`
		Operation           string      `json:"operation"`
		Data                interface{} `json:"data"`
	}

	var opData interface{}

	err := json.Unmarshal([]byte(op.GetDataInJson()), &opData)
	if err != nil {
		return err
	}

	userIdString := ""
	if respUser := op.GetResponsibleUser(); respUser != nil {
		userIdString = respUser.String()
	}

	opToSend := OpToSendStruct{
		ID:                  op.GetID().String(),
		IsRetry:             false,
		Operation:           op.GetOperationType(),
		Data:                opData,
		UserId:              userIdString,
		OperationOriginType: string(opOriginType),
	}

	data, err := json.Marshal(opToSend)
	if err != nil {
		return err
	}

	return l.SecurePublish(data, op.GetChain(), nil, l.config.consumerExchange, priority)
}

func (l *RabbitMQClient) Finish() error {
	err := l.Ch.Close()
	if err != nil {
		return err
	}
	err = l.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (l *RabbitMQClient) GetMsgChan(chain string) <-chan amqp.Delivery {
	l.mu.Lock()
	res := l.Messages[chain]
	l.mu.Unlock()
	return res
}

func (l *RabbitMQClient) SetMsgChan(chain string, c <-chan amqp.Delivery) {
	l.mu.Lock()
	l.Messages[chain] = c
	l.mu.Unlock()
}
