package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	gcfg "gopkg.in/gcfg.v1"

	"github.com/julienschmidt/httprouter"
	"github.com/nararya-anugraha/big-project/db"
	"github.com/nararya-anugraha/big-project/redis"
	"github.com/nararya-anugraha/big-project/user"
	"github.com/nararya-anugraha/big-project/visitor"
)

type ConfigType struct {
	Database db.DatabaseConfigType
	Redis    redis.RedisConfigType
}

func main() {
	os.Exit(Main())
}

func Main() int {
	config := ConfigType{}
	gcfg.ReadFileInto(&config, "config.ini")

	db, err := db.GetDB(&config.Database)
	if err != nil {
		log.Panic(err.Error())
	}

	redisClient := redis.GetRedisClient(&config.Redis)

	router := httprouter.New()

	user.Wire(router, db)
	visitor.Wire(router, redisClient)

	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println("App Started")

	return 0
}
