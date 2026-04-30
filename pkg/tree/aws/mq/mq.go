package mq

type MQ struct {
	Brokers []Broker `tree:"brokers"`
}
