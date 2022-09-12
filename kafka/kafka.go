package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/william094/iris-lib/logx"
	"go.uber.org/zap"
	"strings"
)

type SendMsgCallBackHandler interface {
	OnErrors(producerMsg *sarama.ProducerMessage) error
}

func InitProducer(brokers string, handler SendMsgCallBackHandler) sarama.AsyncProducer {
	address := strings.Split(brokers, ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message
	config.Producer.Return.Errors = true
	var err error
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		logx.SystemLogger.Error("kafka init failed", zap.Error(err))
		panic(err)
	}
	SendCallBack(producer, handler)
	return producer

}

func InitConsumer(brokers, group string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = false
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	address := strings.Split(brokers, ",")
	consumer, err := sarama.NewConsumerGroup(address, group, config)
	if err != nil {
		logx.SystemLogger.Error("kafka consumerGroup init failed", zap.String("group", group), zap.Error(err))
		panic(err)
	}
	go func() {
		for errs := range consumer.Errors() {
			if errs != nil {
				logx.SystemLogger.Error("kafka consume err", zap.Error(err))
			}
		}
	}()
	return consumer
}

func StartKafkaConsumer(group, topics string, client sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler) {
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.SystemLogger.Error("Kafka发送消息异常捕获", zap.Any("err", err))
				return
			}
		}()
		for {
			if err := client.Consume(ctx, strings.Split(topics, ","), handler); err != nil {
				logx.SystemLogger.Error("quit: kafka consumer ", zap.String("group", group), zap.Error(err))
				switch err {
				case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
					// 退出
					logx.SystemLogger.Error("quit: kafka consumer ", zap.String("group", group), zap.Error(err))
					return
				case sarama.ErrOutOfBrokers:
					logx.SystemLogger.Error("kafka 崩溃了~", zap.String("group", group), zap.Error(err))
				default:
					logx.SystemLogger.Error("kafka exception: ", zap.String("group", group), zap.Error(err))
				}
			}
			if err := ctx.Err(); err != nil {
				logx.SystemLogger.Error("Error from context", zap.String("group", group), zap.Error(err))
				break
			}

		}
	}()
}

func SendMsg(topic, key string, data interface{}, producer sarama.AsyncProducer) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.SystemLogger.Error("Kafka发送消息异常捕获", zap.Any("err", err))
				return
			}
		}()
		msg, err := json.Marshal(data)
		if err != nil {
			logx.SystemLogger.Error("Kafka发送消息-消息体序列化失败", zap.String("topic", topic),
				zap.Any("data", data), zap.Error(err))
			return
		}
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(msg),
		}
		logx.SystemLogger.Info("kafka消息发送完成", zap.String("topic", topic), zap.String("key", key),
			zap.Any("data", data))
	}()

}

func SendCallBack(producer sarama.AsyncProducer, handler SendMsgCallBackHandler) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.SystemLogger.Error("发送回调处理异常", zap.Any("err", err))
				return
			}
		}()
		for {
			select {
			case res := <-producer.Successes():
				logx.SystemLogger.Debug("send success", zap.String("topic", res.Topic),
					zap.Int32("partition", res.Partition), zap.Int64("offset", res.Offset))
			case errorMsg := <-producer.Errors():
				keys, _ := errorMsg.Msg.Key.Encode()
				logx.SystemLogger.Error("send failed", zap.String("topic", errorMsg.Msg.Topic),
					zap.String("key", string(keys)), zap.Error(errorMsg.Err))
				if result := handler.OnErrors(errorMsg.Msg); result != nil {
					logx.SystemLogger.Error("发送失败回调处理异常", zap.Error(result))
				}
			}
		}
	}()
}

func Close(producer sarama.AsyncProducer, consumer sarama.ConsumerGroup) {
	if producer != nil {
		producer.Close()
	}
	if consumer != nil {
		consumer.Close()
	}
	logx.SystemLogger.Info("kafka exit。。。。")
}
