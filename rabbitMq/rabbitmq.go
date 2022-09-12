package rabbitMq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/william094/iris-lib/logx"
	"go.uber.org/zap"
	"net"
	"runtime/debug"
	"time"
)

type RabbitMQ struct {
	Url               string
	Vhost             string
	Connection        *amqp.Connection
	Channel           *amqp.Channel
	ReconnectionWatch chan *amqp.Error
}

type Producer struct {
	RabbitMq     *RabbitMQ
	ExchangeName string
	ExchangeType string
	ConfirmWatch chan amqp.Confirmation
}

type Message struct {
	RoutingKey string
	Body       []byte
	Expiration string
	Args       amqp.Table
}

type ConsumerMethod func(msg amqp.Delivery) error

type SendConfirm func(confirms chan amqp.Confirmation, msg *Message)

type Consumer struct {
	RabbitMq     *RabbitMQ
	ExchangeName string
	QueueName    string
	BindingKey   string
	Args         amqp.Table
	Call         ConsumerMethod
}

func MakeMessage(routingKey, expiration string, body []byte, args amqp.Table) *Message {
	return &Message{
		RoutingKey: routingKey,
		Body:       body,
		Expiration: expiration,
		Args:       args,
	}
}

func OpenConnection(vhost, url string) *RabbitMQ {
	rabbitmq, err := Connection(vhost, url)
	if err != nil {
		panic(err)
	}
	return rabbitmq
}

func Connection(vhost, url string) (*RabbitMQ, error) {
	mqConfig := amqp.Config{
		Heartbeat: time.Second * 5,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 20*time.Minute)
		},
	}
	if vhost != "" {
		mqConfig.Vhost = vhost
	}
	conn, err := amqp.DialConfig(url, mqConfig)
	if err != nil {
		logx.SystemLogger.Info("MQ连接失败...")
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		logx.SystemLogger.Info("MQ通道打开失败...")
		return nil, err
	}
	return &RabbitMQ{
		Url:        url,
		Vhost:      vhost,
		Connection: conn,
		Channel:    ch,
	}, nil
}

func Reconnection(rabbitmq *RabbitMQ) {
	for {
		select {
		case s1 := <-rabbitmq.Connection.NotifyClose(make(chan *amqp.Error, 1)):
			logx.SystemLogger.Error("connection notify close", zap.Error(s1))
			if r, err := Connection(rabbitmq.Vhost, rabbitmq.Url); err != nil {
				logx.SystemLogger.Error("connection notify close reconnection error", zap.Error(err))
			} else {
				rabbitmq.Connection = r.Connection
				rabbitmq.Channel = r.Channel
				logx.SystemLogger.Error("connection notify close,reconnect success")
			}
		case s2 := <-rabbitmq.Channel.NotifyClose(make(chan *amqp.Error, 1)):
			logx.SystemLogger.Error("channel notify close", zap.Error(s2))

			if r, err := Connection(rabbitmq.Vhost, rabbitmq.Url); err != nil {
				logx.SystemLogger.Error("channel notify close reconnection failed", zap.Error(err))
			} else {
				rabbitmq.Connection = r.Connection
				rabbitmq.Channel = r.Channel
				logx.SystemLogger.Error("channel notify close,reconnect success")
			}
		}
	}
}

func NewProducer(exchangeName, exchangeType string, rabbitmq *RabbitMQ) *Producer {
	producer := &Producer{
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	producer.RabbitMq = rabbitmq
	//声明交换器
	err := producer.RabbitMq.Channel.ExchangeDeclare(
		exchangeName, //交换机名称
		exchangeType, //交换机类型
		true,         //是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能
		false,        //没有剩余绑定时将在服务器重新启动之前和之后被删除。
		false,        //true:声明为内部交换器，不接收消费发布
		false,        //true：无需等待服务器确认即可声明。通道可能因错误而关闭。添加 NotifyClose 侦听器以响应任何异常
		nil,          //其他参数
	)
	if err != nil {
		logx.SystemLogger.Info("MQ通道声明交换器失败...", zap.String("exchangeName", exchangeName))
		//panic(err)
	}
	if err := producer.RabbitMq.Channel.Confirm(false); err != nil {
		logx.SystemLogger.Info("MQ通道发送确认开启失败...", zap.String("exchangeName", exchangeName))
		//panic(err)
	}
	//声明一个发送确认回调chan
	producer.ConfirmWatch = producer.RabbitMq.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	go Reconnection(rabbitmq)
	return producer
}

func StartRabbitConsumer(exchangeName, queueName, bindingKey string, args amqp.Table, method ConsumerMethod, rabbitmq *RabbitMQ) {
	consumer := &Consumer{
		ExchangeName: exchangeName,
		QueueName:    queueName,
		BindingKey:   bindingKey,
		Args:         args,
		Call:         method,
	}
	consumer.RabbitMq = rabbitmq
	//声明一个队列
	consumer.RabbitMq.Channel.QueueDeclare(
		queueName, //队列名称
		true,      //是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能
		false,     //true:没有活动的消费者将被删除
		false,     //true 表示只能由声明它们的连接访问，并且在连接关闭时将被删除。尝试声明、绑定、消费、清除或删除同名队列时，其他连接上的通道将收到错误
		true,      //当 noWait 为 true 时，队列将假定在服务器上声明。如果有满足条件的队列或尝试从不同的连接修改现有队列，则会出现通道异常
		nil,       //其他参数
	)
	//交换器队列绑定
	consumer.RabbitMq.Channel.QueueBind(
		queueName,    //队列名称
		bindingKey,   //绑定路由key
		exchangeName, //交换器名称
		false,        //当 noWait 为 false 且无法绑定队列时，通道将因错误而关闭
		args,         //其他参数
	)
	delivery, err := consumer.RabbitMq.Channel.Consume(
		queueName, //监听队列名称
		"",        //监听消费服务名称
		false,     //是否自动提交ACK
		false,     //当exclusive 为true 时，服务器将确保这是该队列中的唯一消费者。 当exclusive 为false 时，服务器将在多个消费者之间公平地分发交付
		false,     //使用临时队列接收
		false,     //当 noWait 为 true 时，不要等待服务器确认请求并立即开始交付。如果无法消费，则会引发通道异常并关闭通道
		nil)
	if err != nil {
		logx.SystemLogger.Info("消费监听失败", zap.Error(err), zap.String("strack", string(debug.Stack())))
		panic(err)
	}
	go Reconnection(rabbitmq)
	for {
		select {
		case msg := <-delivery:
			if rabbitmq.Connection.IsClosed() || rabbitmq.Channel.IsClosed() {
				msg.Reject(true)
				break
			}
			if result := method(msg); result != nil {
				//消费失败
				msg.Nack(false, true)
			} else {
				//消费成功
				msg.Ack(true)
			}
		}
	}
}

func (p *Producer) SendMessage(msg *Message, confirmMethod SendConfirm) error {
	if confirmMethod != nil {
		defer confirmMethod(p.ConfirmWatch, msg)
	} else {
		defer confirmOne(p.ConfirmWatch, msg)
	}
	if err := p.RabbitMq.Channel.Publish(p.ExchangeName, msg.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msg.Body,
		Headers:      msg.Args,
		DeliveryMode: 2,              //2:消息持久化到磁盘
		Expiration:   msg.Expiration, //过期时间
	}); err != nil {
		logx.SystemLogger.Info("消息发送失败", zap.Any("messageInfo", msg))
		return err
	}
	logx.SystemLogger.Info("消息发送成功", zap.Any("messageInfo", msg))
	return nil
}

func (r *RabbitMQ) Close() {
	if err := r.Channel.Close(); err != nil {
		logx.SystemLogger.Info("channel close failed")
	}
	if err := r.Connection.Close(); err != nil {
		logx.SystemLogger.Info("connection close failed")
	}
	logx.SystemLogger.Info("close ")
}

// 消息确认
func confirmOne(confirms chan amqp.Confirmation, msg *Message) {
	if confirmed := <-confirms; confirmed.Ack {
		logx.SystemLogger.Debug("MQ发送确认成功", zap.Any("body", msg), zap.Any("delivery tag", confirmed.DeliveryTag))
	} else {
		//可以尝试重发
		logx.SystemLogger.Info("MQ发送确认失败", zap.Any("body", msg), zap.Any("delivery tag", confirmed.DeliveryTag))
	}
}
