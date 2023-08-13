package rabbitmqClient

import (
	"fmt"
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
	}
}

func NewRabbitMQClient(user, pass, host, port, consumerExchange string, consumerQueues map[string]string) (*RabbitMQClient, error) {
	l := &RabbitMQClient{config: &struct {
		user             string
		pass             string
		host             string
		port             string
		consumerQueues   map[string]string
		consumerExchange string
	}{
		user:             user,
		pass:             pass,
		host:             host,
		port:             port,
		consumerQueues:   consumerQueues,
		consumerExchange: consumerExchange,
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
