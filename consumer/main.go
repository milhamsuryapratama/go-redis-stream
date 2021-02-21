package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/milhamsuryapratama/go-redis-stream/consumer/handler"
	"github.com/milhamsuryapratama/go-redis-stream/packages/event"
	"github.com/milhamsuryapratama/go-redis-stream/utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var (
	client        *redis.Client
	waitGrp       sync.WaitGroup
	streamName    string = "ilhamStream"
	consumerGroup string = "ilhamGroup"
)

func init() {
	var err error
	client, err = utils.NewRedisClient()
	if err != nil {
		panic(err)
	}

	createConsumerGroup()
}

func main() {
	go consumeEvent()

	chanOS := make(chan os.Signal)
	//Gracefully disconection
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)
	<-chanOS

	waitGrp.Wait()
	client.Close()
}

func createConsumerGroup() {
	if _, err := client.XGroupCreateMkStream(context.Background(), streamName, consumerGroup, "0").Result(); err != nil {
		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)
			panic(err)
		}
	}
}

func consumeEvent() {
	for {
		func() {
			streams, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Streams:  []string{streamName, ">"},
				Group:    consumerGroup,
				Consumer: "Ilham",
				Count:    10,
				Block:    0,
			}).Result()

			if err != nil {
				log.Printf("err on consume events: %+v\n", err)
				return
			}

			for _, stream := range streams[0].Messages {
				waitGrp.Add(1)
				// process stream
				go processStream(stream, handler.HandlerFactory())
			}
			waitGrp.Wait()
		}()
	}
}

func processStream(stream redis.XMessage, handlerFactory func(event event.Type) handler.Handler) {

	defer waitGrp.Done()

	typeEvent := stream.Values["type"].(string)
	newEvent, _ := event.New(event.Type(typeEvent))

	err := newEvent.UnmarshalBinary([]byte(stream.Values["data"].(string)))
	if err != nil {
		fmt.Printf("error on unmarshal stream:%v\n", stream.ID)
		return
	}

	newEvent.SetID(stream.ID)

	h := handlerFactory(newEvent.GetType())
	err = h.Handle(newEvent)

	if err != nil {
		fmt.Printf("error on process event:%v\n", newEvent)
		fmt.Println(err)
		return
	}

	client.XAck(context.Background(), streamName, consumerGroup, stream.ID)
}
