package nsq

import (
	"github.com/nsqio/go-nsq"
	"go-central-control/mq"
	"go-central-control/utils"
)

type NsqCustomer mq.Consumer

func (c NsqCustomer) StartConsumer(topic string, channel string, callback func()) {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		utils.Log.Fatal(err)
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		utils.Log.Info(c, string(message.Body))
		return nil
	}))
	// 连接到单例nsqd
	if err := consumer.ConnectToNSQD(c.Ip + ":" + c.Port); err != nil {
		utils.Log.Fatal(err)
	}
	<-consumer.StopChan
}
