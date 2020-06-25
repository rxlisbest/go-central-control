package mq

type ConsumerInterface interface {
	StartConsumer(topic string, channel string, callback func())
}

type Consumer struct {
	Ip   string
	Port string
}
