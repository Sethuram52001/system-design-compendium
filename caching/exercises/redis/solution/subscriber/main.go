package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	sub := redisClient.Subscribe("user.events")
	defer sub.Close()

	for msg := range sub.Channel() {
		log.Printf("received event: %s", msg.Payload)
	}
}
