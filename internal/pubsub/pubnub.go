package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type PubNub struct {
	pn *pubnub.PubNub
}

func NewPubSub(pn *pubnub.PubNub) *PubNub {
	return &PubNub{
		pn: pn,
	}
}

func (p *PubNub) Subscribe(channel string, dbClient *mongo.Client) error {
	listener := pubnub.NewListener()
	p.pn.AddListener(listener)

	p.pn.Subscribe().
		Channels([]string{channel}).
		Execute()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					log.Printf("Subscribing to channel '%v'... connected", channel)
				case pubnub.PNBadRequestCategory:
					log.Printf("Subscribing to channel '%v'... bad request", channel)
				default:
					log.Println(status)
				}
			case message := <-listener.Message:
				var requestBody map[string]interface{}
				request, _ := json.Marshal(&message.Message)
				json.Unmarshal(request, &requestBody)
				fmt.Println(requestBody)
				collection := dbClient.Database(viper.GetString("mongodb.database")).Collection(channel)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				res, err := collection.InsertOne(ctx, requestBody)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(res)
				}
			}
		}
	}()
	return nil
}
