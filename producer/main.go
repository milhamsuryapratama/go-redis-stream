package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/milhamsuryapratama/go-redis-stream/packages/event"
	"github.com/milhamsuryapratama/go-redis-stream/utils"
	"math/rand"
	"time"
)

const (
	MaxUserIDRange = 10000
)

var (
	client *redis.Client
)

func init() {
	var err error
	client, err = utils.NewRedisClient()
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		amount int = 1000
	)

	eventType := []event.Type{event.TrxType}[rand.Intn(1)]

	for i := 0; i < 5; i++ {
		userID := int(rand.Intn(MaxUserIDRange))

		newID, _ := produceMsg(map[string]interface{}{
			"type": string(eventType),
			"data": &event.TransactionEvent{
				Base: &event.Base{
					Type:     eventType,
					DateTime: time.Now(),
				},
				UserID: userID,
				Amount: userID * amount,
			},
		})

		fmt.Println("NewID", newID)
	}
}

func produceMsg(event map[string]interface{}) (string, error) {
	ev, _ := json.MarshalIndent(event, "", "  ")
	fmt.Println(string(ev))

	return client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "ilhamStream",
		Values: event,
	}).Result()
}
