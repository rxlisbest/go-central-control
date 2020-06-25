package main

import (
	"go-central-control/mq"
	"log"
	"math/rand"
	"strconv"
	"time"
	"github.com/nsqio/go-nsq"
	"go-central-control/mq"
	"go-central-control/mq/nsq"
)

func main() {
	var consumer mq.ConsumerInterface
	consumer = nsq.NsqCustomer{"127.0.0.1" , "4150"}
	go startConsumer(1)
	go startConsumer(2)
	go startConsumer(3)
	startProducer()
}

// 生产者
func startProducer() {
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 发布消息
	for {
		if err := producer.Publish("test", []byte(strconv.Itoa(rand.Intn(100)))); err != nil {
			log.Fatal("publish error: " + err.Error())
		}

		time.Sleep(100 * time.Second)
	}
}
