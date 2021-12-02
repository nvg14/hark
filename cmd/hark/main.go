package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nvg14/hark/internal/cron"
	"github.com/nvg14/hark/internal/database"
	"github.com/nvg14/hark/internal/pubsub"
	pubnub "github.com/pubnub/go"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type Agents []struct {
	PublishKey   string
	SubscribeKey string
	SecretKey    string
	UUID         string
	Channels     []string
}

var (
	app        *cli.App
	configPath string
)

func main() {
	exit := make(chan string)
	app = cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Required:    false,
			Destination: &configPath,
			Usage:       "Path to a configuration `FILE`",
		},
	}
	database.DirectoryPath = "./PubnubMessages/"
	if configPath == "" {
		configPath = "./config"
	}

	err := setViper()
	if err != nil {
		log.Println(err)
		return
	}

	dbClient, err := database.NewDatabase()
	if err != nil {
		log.Println(err)
		return
	}

	folders:= ChannelsToFolders()
	interval:= viper.GetInt("s3.interval")
	bucket:= viper.GetString("s3.bucket")

	app.Action = func(c *cli.Context) error {
		setAgents(dbClient)
		go cron.RunUploadS3(time.Now().Hour(), time.Now().Minute(),time.Now().Second(),interval,bucket,dbClient,folders)
		return nil
	}

	// Run the CLI app
	err = app.Run(os.Args)
	if err != nil {
		fmt.Print(err.Error())
	}

	for {
		select {
		case <-exit:
			os.Exit(0)
		}
	}
}

func setViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func setAgents(dbClient *s3manager.Uploader) error {
	var tempAgents Agents
	viper.UnmarshalKey("agents", &tempAgents)
	for _, agent := range tempAgents {
		log.Printf("PublishKey: %v | SubscribeKey: %v | SecretKey: %v | UUID: %v", agent.PublishKey, agent.SubscribeKey, agent.SecretKey, agent.UUID)
		log.Printf("Listening channel: %s", agent.Channels)
		config := pubnub.NewConfig()
		config.PublishKey = agent.PublishKey
		config.SubscribeKey = agent.SubscribeKey
		config.UUID = agent.UUID
		config.SecretKey = agent.SecretKey
		pn := pubnub.NewPubNub(config)
		ps := pubsub.NewPubSub(pn)
		fmt.Println("subscribing to channel :",agent.Channels)
		err := ps.Subscribe(agent.Channels)
		if err != nil {
			return err
		}
	}
	return nil
}


func ChannelsToFolders() []string{
	var Agents Agents
	viper.UnmarshalKey("agents", &Agents)
	FilePaths := []string{}
	for _,agent:=range Agents {
		for _,channel := range agent.Channels{
			s := strings.Split(channel, "|")
			directoryPath := s[1] + "/" + s[2] + "/" + s[3]
			FilePaths = append(FilePaths, directoryPath)
		}
	}
	return FilePaths
}