package mq

type ProducerInterface interface {
	StartProducer(topic string, callback func())
}

type Producer struct {
	Ip   string
	Port string
}