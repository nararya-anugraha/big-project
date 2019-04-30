package visitor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// visitor.handler

type visitorModuleType struct {
	redisClient *redis.Client
}

func Wire(router *httprouter.Router, redisClient *redis.Client) {
	visitorModule := visitorModuleType{
		redisClient: redisClient,
	}
	router.GET("/api/visitor", visitorModule.getVisitorHandler)
}

type visitorResponseType struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	VisitorCount int    `json:"visitor_count"`
}

func (visitorModule *visitorModuleType) getVisitorHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := "key:big-project:visitor-count"

	visitorCountString, err := visitorModule.redisClient.Get(key).Result()
	if err == redis.Nil {
		//Key doesn't exist
		visitorCountString = "0"
	}

	visitorCount, err := strconv.Atoi(visitorCountString)
	if err != nil {
		handleError(writer, err)
		return
	}

	visitorCount++
	err = visitorModule.redisClient.Set(key, visitorCount, 0).Err()
	if err != nil {
		handleError(writer, err)
		return
	}

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(visitorCount)
	if err != nil {
		handleError(writer, err)
		return
	}
}

func handleError(writer http.ResponseWriter, err error) {
	encoder := json.NewEncoder(writer)
	encoder.Encode(err)

	writer.Header().Add("status", "500")
	writer.Header().Add("content-type", "application/json")

}
