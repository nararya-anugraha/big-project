package nsq

import (
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
)

type NSQConfigType struct {
	MaxAttempts           uint16
	MaxInFlight           int
	Topic                 string
	Channel               string
	ConsumerLookupAddress string
	ProducerNSQAddress    string
}

type NSQConsumerType struct {
	consumer              *nsq.Consumer
	consumerLookupAddress string
}

type NSQProducerType struct {
	topic    string
	producer *nsq.Producer
}

func CreateConsumer(config *NSQConfigType) (*NSQConsumerType, error) {
	nsqConf := nsq.NewConfig()
	nsqConf.MaxAttempts = config.MaxAttempts
	nsqConf.MaxInFlight = config.MaxInFlight

	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &NSQConsumerType{
		consumer:              consumer,
		consumerLookupAddress: config.ConsumerLookupAddress,
	}, nil
}

func (nsqConsumer *NSQConsumerType) AddHandler(handler nsq.HandlerFunc) {
	nsqConsumer.consumer.AddHandler(handler)
}

func (nsqConsumer *NSQConsumerType) Run() {
	err := nsqConsumer.consumer.ConnectToNSQLookupd(nsqConsumer.consumerLookupAddress)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateProducer(config *NSQConfigType) (*NSQProducerType, error) {
	producer, err := nsq.NewProducer(config.ProducerNSQAddress, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	return &NSQProducerType{
		producer: producer,
		topic:    config.Topic,
	}, nil
}

func (nsqProducer *NSQProducerType) Publish(message string) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return nsqProducer.producer.Publish(nsqProducer.topic, payload)
}
