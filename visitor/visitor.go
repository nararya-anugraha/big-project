package visitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	nsqclient "github.com/nararya-anugraha/big-project/nsq"
	"github.com/nsqio/go-nsq"
)

// visitor.handler

const key = "key:big-project:visitor-count"

type visitorModuleType struct {
	redisClient *redis.Client
	consumer    *nsqclient.NSQConsumerType
	producer    *nsqclient.NSQProducerType
}

func Wire(router *httprouter.Router, redisClient *redis.Client, consumer *nsqclient.NSQConsumerType, producer *nsqclient.NSQProducerType) {
	visitorModule := visitorModuleType{
		redisClient: redisClient,
		consumer:    consumer,
		producer:    producer,
	}
	router.GET("/api/visitor", visitorModule.getVisitorHandler)
	consumer.AddHandler(visitorModule.incrementVisitorCount)
}

type visitorResponseType struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	VisitorCount int    `json:"visitor_count"`
}

func (visitorModule *visitorModuleType) getVisitorHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	visitorCountString, err := visitorModule.redisClient.Get(key).Result()
	if err == redis.Nil {
		//Key doesn't exist
		visitorCountString = "1"
	}

	visitorCount, err := strconv.Atoi(visitorCountString)
	if err != nil {
		handleError(writer, err)
		return
	}

	visitorModule.producer.Publish("incrementVisitorCount")

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(visitorCount)
	if err != nil {
		handleError(writer, err)
		return
	}
}

func (visitorModule *visitorModuleType) incrementVisitorCount(message *nsq.Message) error {
	var messageString string

	json.Unmarshal(message.Body, &messageString)
	fmt.Println(messageString)
	if messageString != "incrementVisitorCount" {
		return nil
	}

	err := visitorModule.redisClient.Incr(key).Err()

	if err != nil {
		return err
	}

	return nil

}

func handleError(writer http.ResponseWriter, err error) {
	encoder := json.NewEncoder(writer)
	encoder.Encode(err)
	writer.Header().Add("status", "500")
	writer.Header().Add("content-type", "application/json")
}
