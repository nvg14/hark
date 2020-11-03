package pubsub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	pubnub "github.com/pubnub/go"
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
				guid, err := uuid.NewRandom()
				if err != nil {
					log.Println(err)
				}
				SaveMetadata(requestBody, "/tmp/hark/"+guid.String())
				// collection := dbClient.Database(viper.GetString("mongodb.database")).Collection(channel)
				// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				// defer cancel()
				// res, err := collection.InsertOne(ctx, requestBody)
				// if err != nil {
				// 	log.Println(err)
				// } else {
				// 	log.Println(res)
				// }
			}
		}
	}()
	return nil
}

func SaveMetadata(meta interface{}, outputFilepath string) error {
	data, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return err
	}
	err = CreatePathIfNotExist(outputFilepath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFilepath, []byte(data), 0766)
	return err
}

func CreatePathIfNotExist(filePath string) error {
	path, _ := filepath.Abs(filePath)
	parentDir := filepath.Dir(path)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		err = os.MkdirAll(parentDir, 0766)
		if err != nil {
			log.Println("Cant create directory")
			return err
		}
	}
	return nil
}
