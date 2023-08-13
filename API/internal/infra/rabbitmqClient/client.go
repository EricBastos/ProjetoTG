package rabbitmqClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQClientConfig struct {
	User             string
	Pass             string
	Host             string
	Port             string
	ProducerExchange string
}

type DeliveryMsg amqp.Delivery

type RabbitMQClient struct {
	ch     *amqp.Channel
	conn   *amqp.Connection
	config *RabbitMQClientConfig
}

func NewRabbitMQClient(
	config *RabbitMQClientConfig,
) (*RabbitMQClient, error) {
	result := &RabbitMQClient{}
	result.config = config
	err := result.initialize()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l *RabbitMQClient) initialize() error {
	rabbitmqConnString := fmt.Sprintf("amqp://%s:%s@%s:%s/vhost",
		l.config.User,
		l.config.Pass,
		l.config.Host,
		l.config.Port,
	)
	var err error
	l.conn, err = amqp.Dial(rabbitmqConnString)
	if err != nil {
		return err
	}

	l.ch, err = l.conn.Channel()
	if err != nil {
		return err
	}

	err = l.ch.Confirm(false)
	if err != nil {
		return err
	}

	err = l.ch.ExchangeDeclare(
		l.config.ProducerExchange,
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

	go l.observeConnection()

	return nil

}

func (l *RabbitMQClient) observeConnection() {
	connErr := <-l.conn.NotifyClose(make(chan *amqp.Error, 2))
	log.Println(fmt.Sprintf("(RabbitMQ) Lost rabbitmq connection due to error: %s", connErr))
	for {
		l.closeConnections()
		if err := l.initialize(); err == nil {
			log.Println("(RabbitMQ) Connection to rabbitmq reestablished")
			return
		}
		time.Sleep(30 * time.Second)
	}
}

func (l *RabbitMQClient) closeConnections() {
	if l.ch != nil {
		_ = l.ch.Close()
	}
	if l.conn != nil && l.conn.IsClosed() {
		_ = l.conn.Close()
	}
}

func (l *RabbitMQClient) CallSmartcontract(op entities.SmartContractOp, opOriginType entities.OperationOriginType) error {

	type OpToSendStruct struct {
		ID                  string      `json:"id"`
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
		Operation:           op.GetOperationType(),
		Data:                opData,
		UserId:              userIdString,
		OperationOriginType: string(opOriginType),
	}

	data, err := json.Marshal(opToSend)
	if err != nil {
		return err
	}

	return l.SecurePublish(data, nil, l.config.ProducerExchange, op.GetChain())
}

func (l *RabbitMQClient) SecurePublish(data []byte, headers map[string]interface{}, exchange string, topic string) error {
	conf, err := l.ch.PublishWithDeferredConfirmWithContext(
		context.Background(),
		exchange, // exchange
		topic,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
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

func (l *RabbitMQClient) Ack(tag uint64, multiple bool) error {
	return l.ch.Ack(tag, multiple)
}

func (l *RabbitMQClient) Nack(tag uint64, multiple bool, requeue bool) error {
	return l.ch.Nack(tag, multiple, requeue)
}

func (l *RabbitMQClient) Finish() error {
	err := l.ch.Close()
	if err != nil {
		return err
	}
	err = l.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
