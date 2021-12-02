package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pubnub "github.com/pubnub/go"
)

type PubNub struct {
	pn *pubnub.PubNub
}

func NewPubSub(pn *pubnub.PubNub) *PubNub {
	return &PubNub{
		pn: pn,
	}
}

func (p *PubNub) Subscribe(channels []string) error {
	listener := pubnub.NewListener()
	p.pn.AddListener(listener)

	p.pn.Subscribe().
		Channels(channels).
		Execute()

	go func() {
		//	fmt.Println(channel)
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					log.Printf("Subscribing to channel '%v'... connected", channels)
				case pubnub.PNBadRequestCategory:
					log.Printf("Subscribing to channel '%v'... bad request", channels)
				default:
					log.Println(status)
				}
			case message := <-listener.Message:
				fmt.Println(message.Message)
				go ProcessMessage(message)

			}
		}
	}()
	return nil
}

func ProcessMessage(message *pubnub.PNMessage){
	timeStamp := time.Now()
	date := timeStamp.Format("2006-01-02_15")

	directoryPath := "./PubnubMessages"

	s := strings.Split(message.Channel, "|")

	//create directory
	directoryPath = directoryPath + "/" + s[1] + "/" + s[2] + "/" + s[3]
	if err := ensureDir(directoryPath); err != nil {
		fmt.Println("Failed to create directory")
		os.Exit(1)

	}

	//file path
	path := directoryPath + "/" + date + ".txt"
	//path := "logs.txt"
	//check file if exists
	var _, err = os.Stat(path)


	//write message to file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//var requestBody map[string]interface{}
	request, _ := json.Marshal(&message)
	if _, err := f.Write(request); err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("processed the message")
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
