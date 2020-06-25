package main

import (
	"go-central-control/mq"
	"go-central-control/mq/nsq"
)

func main() {
	var consumer mq.ConsumerInterface
	consumer = nsq.NsqCustomer{"127.0.0.1", "4150"}
	go func() { consumer.StartConsumer("test", "1", func() {

	}) }()

	var producer mq.ProducerInterface
	producer = nsq.NsqProducer{"127.0.0.1", "4150"}
	producer.StartProducer("test", func() {

	})
}
