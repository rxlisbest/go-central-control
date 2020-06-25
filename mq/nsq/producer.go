package nsq

import (
	"github.com/nsqio/go-nsq"
	"go-central-control/mq"
	"go-central-control/utils"
	"math/rand"
	"strconv"
	"time"
)

type NsqProducer mq.Producer

func (p NsqProducer) StartProducer(topic string, callback func()) {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(p.Ip + ":" + p.Port, cfg)
	if err != nil {
		utils.Log.Fatal(err)
	}

	// 发布消息
	for {
		if err := producer.Publish(topic, []byte(strconv.Itoa(rand.Intn(100)))); err != nil {
			utils.Log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}
