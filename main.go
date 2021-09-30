package main

import (
	"alfred/bot"
	"alfred/config"
	"alfred/registry"
	"context"

	"github.com/Strum355/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
)

func main() {
	// Load the configuration for bot
	config.Load()

	// Setup the logger
	log.InitSimpleLogger(&log.Config{})

	// Print out the loaded config
	config.PrintSettings()

	// Create a service registry to keep track of the connected services
	serviceReg := registry.NewRegistry()
	// serviceReg.Services = append(serviceReg.Services, registry.Service{
	// 	Name:           "Sample Service",
	// 	CommandTrigger: "sample",
	// 	URL:            "127.0.0.1",
	// 	Port:           "3000",
	// 	Alive:          true,
	// })

	serviceReg.PrintRegistry()

	// // Connect to the Slack API using socket mode to be able to receive messages
	client := slack.New(viper.GetString("alfred.token"), slack.OptionAppLevelToken(viper.GetString("alfred.apptoken")))

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
	)

	bot := bot.NewBot(client, socketClient, serviceReg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go bot.ListenForWebhook(ctx)

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Info("Alfred shutting down")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Error("Could not cast event")
						continue
					}
					socketClient.Ack(*event.Request)
					err := bot.HandleEvent(ctx, eventsAPIEvent)
					if err != nil {
						log.WithError(err).Error("Could not handle event")
						continue
					}
				}
			}
		}
	}(ctx, client, socketClient)

	// Run the socket client
	socketClient.Run()
}
