package rabbitmq

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	Client *amqp.Connection
	Queue  map[string]*amqp.Queue
)

// NewRabbitMQ
// @author: lipper
// @function: NewRabbitMQ
// @description: 建立rabbitmq 连接
// @return: *amqp.Connection, map[string]*amqp.Queue
func NewRabbitMQ(ctx context.Context, conf MQConf, queues ...QueueDeclare) (*amqp.Connection, map[string]*amqp.Queue) {
	var (
		err     error
		address string
		ch      *amqp.Channel
	)

	Queue = make(map[string]*amqp.Queue)

	address = fmt.Sprintf("amqp://%s:%s@%s:%d/",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
	)

	logrus.Debugf("rabbitmq address: %s", address)

	if Client, err = amqp.DialConfig(address, amqp.Config{
		Vhost:           conf.Vhost,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}); err != nil {
		logrus.Panicf("dail rabbitmq[%s] err: %v", address, err)
	}

	if ch, err = Client.Channel(); err != nil {
		logrus.Panicf("rabbitmq client connection[%s] get new channel err: %v", address, err)
	}

	go func() {
		select {
		case <-ctx.Done():
			logrus.Warn("rabbitmq client connection closed by gctx")
			_ = Client.Close()
		}
	}()

	if len(queues) > 0 {
		for _, queue := range queues {
			var (
				dq amqp.Queue
			)

			if dq, err = ch.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, false, false, nil); err != nil {
				logrus.Panicf("rabbitmq client[%s] declare queue[%s] err: %v", address, queue.Name, err)
			}

			if err = ch.QueueBind(queue.Name, conf.Vhost, queue.Exchange, false, nil); err != nil {
				logrus.Panicf("rabbitmq bind queue[%s] to vhost[%s] exchange[%s] err: %v", queue.Name, conf.Vhost, queue.Exchange, err)
			}

			Queue[dq.Name] = &dq
		}
	}

	logrus.Infof("rabbitmq client[%s] init success!!!", address)

	return Client, Queue
}
