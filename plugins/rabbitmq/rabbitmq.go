package rabbitmq

import (
	"fmt"
	"log"

	"github.com/gookit/color"
	"github.com/streadway/amqp"
)

type Consumer func(amqp.Delivery)

// 消息体
type Message struct {
	DelayTime int
	Body      string
}

type MessageQueue struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// 消费者回调方法
type CallBack func(amqp.Delivery)

// @author: lipper
// @object: *MessageQueue
// @function: NewRabbitMQ
// @description: 建立rabbitmq 连接
// @return: MessageQueue MessageQueue , err error
func NewRabbitMQ(conf MQConf) (MessageQueue, error) {
	messageQueue := MessageQueue{}
	// 建立amqp链接
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Vhost),
	)
	if err != nil {
		color.Redln("rabbitmq 连接失败!" + err.Error())
		return messageQueue, err
	}
	messageQueue.Conn = conn

	// 建立channel通道
	ch, err := conn.Channel()
	if err != nil {
		color.Redln("rabbitmq 建立channel失败!" + err.Error())
		return messageQueue, err
	}
	messageQueue.Ch = ch
	return messageQueue, err
}

// @author: lipper
// @object: *MessageQueue
// @function: Close
// @description: 关闭链接
func (mq *MessageQueue) Close() {
	mq.Ch.Close()
	mq.Conn.Close()
}

// @author: lipper
// @object: *MessageQueue
// @function: DeclareQueue
// @description: 创建队列
// @param: queueName string, durable, autoDelete, exclusive, nowait bool, args amqp.Table
// @return: amqp.Queue, error
func (mq *MessageQueue) DeclareQueue(queueName string, durable, autoDelete, exclusive, nowait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := mq.Ch.QueueDeclare(queueName, durable, autoDelete, exclusive, nowait, args)
	return q, err
}

// @author: lipper
// @object: *MessageQueue
// @function: DeclareExchange
// @description: 创建交换器
// @param: exchaneName, exchangeType string, durable, autoDelete, internal, nowait bool, args amqp.Table
// @return:  error
func (mq *MessageQueue) DeclareExchange(exchaneName, exchangeType string, durable, autoDelete, internal, nowait bool, args amqp.Table) error {
	err := mq.Ch.ExchangeDeclare(exchaneName, exchangeType, durable, autoDelete, internal, nowait, args)
	return err
}

// @author: lipper
// @object: *MessageQueue
// @function: BindQueue
// @description: 绑定队列
// @param: queueName, routeKey, exchangeName string
// @return:  error
func (mq *MessageQueue) BindQueue(queueName, routeKey, exchangeName string) error {
	err := mq.Ch.QueueBind(queueName, routeKey, exchangeName, false, nil)
	return err
}

// @author: lipper
// @object: *MessageQueue
// @function: ConsumeMessage
// @description: 消费消息
// @param: callBack CallBack, queueName, consumer string, noAck, exclusive, noLocal, noWait bool
// @return:
func (mq *MessageQueue) ConsumeMessage(callBack CallBack, queueName, consumer string, noAck, exclusive, noLocal, noWait bool) {
	// 设置Qos
	err := mq.Ch.Qos(1, 0, false)
	if err != nil {
		color.Redln("设置Qos失败!" + err.Error())
		return
	}
	// 监听消息
	msgs, err := mq.Ch.Consume(queueName, consumer, noAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return
	}
	var forever chan struct{}
	go func() {
		for d := range msgs {
			callBack(d)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// @author: lipper
// @object: *MessageQueue
// @function: ConsumeMessage
// @description: 生产消息
// @param: message string
// @return: error
func (mq *MessageQueue) SendMessage(exchangeName, routingKey, message string) error {
	defer mq.Conn.Close()
	err := mq.Ch.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	},
	)
	return err
}
