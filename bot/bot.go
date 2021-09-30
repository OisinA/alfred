package bot

import (
	"alfred/registry"
	"context"
	"strings"

	"github.com/Strum355/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Bot struct {
	client       *slack.Client
	socketClient *socketmode.Client
	registry     registry.Registry
}

func NewBot(client *slack.Client, socketClient *socketmode.Client, serviceReg registry.Registry) Bot {
	return Bot{
		client:       client,
		socketClient: socketClient,
		registry:     serviceReg,
	}
}

func (b *Bot) HandleEvent(ctx context.Context, event slackevents.EventsAPIEvent) error {
	// Switch on event type
	switch event.Type {
	// Access the inner event of a callback
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// Event is when someone @ mentions the bot
			if ev.Text == "" {
				return nil
			}
			log.WithFields(log.Fields{
				"event": ev,
			}).Info("Mentioned")

			// Remove the mention and extract the command
			withoutMention := strings.Join(strings.Split(ev.Text, " ")[1:], " ")
			command := strings.Split(withoutMention, " ")[0]
			withoutCommand := strings.Join(strings.Split(withoutMention, " ")[1:], " ")

			// Fetch the user info
			user, err := b.client.GetUserInfo(ev.User)
			if err != nil {
				return err
			}

			// Send the command to the relevant service
			msg, err := b.registry.SendByCommand(command, registry.SendCommand{
				Command: command,
				User:    user.Name,
				Args:    withoutCommand,
			})

			// If there is an error, print it to the user
			if err != nil {
				attachment := slack.Attachment{
					Pretext: "Error",
					Text:    err.Error(),
				}
				_, _, err := b.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
				if err != nil {
					return err
				}
				return nil
			}

			attachment := slack.Attachment{
				Pretext: "Response",
				Text:    msg,
			}
			_, _, err = b.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
			if err != nil {
				return err
			}

			return nil
		case *slackevents.MessageEvent:
			if ev.Text == "" || ev.ChannelType == "channel" {
				return nil
			}
			log.WithFields(log.Fields{
				"message": ev.Text,
			}).Info("Received message")
			command := strings.Split(ev.Text, " ")[0]
			withoutCommand := strings.Join(strings.Split(ev.Text, " ")[1:], " ")

			user, err := b.client.GetUserInfo(ev.User)
			if err != nil {
				return err
			}
			msg, err := b.registry.SendByCommand(command, registry.SendCommand{
				Command: command,
				User:    user.Name,
				Args:    withoutCommand,
			})

			if err != nil {
				attachment := slack.Attachment{
					Pretext: "Error",
					Text:    err.Error(),
				}
				_, _, err := b.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
				if err != nil {
					return err
				}

				return nil
			}

			attachment := slack.Attachment{
				Pretext: "Response",
				Text:    msg,
			}
			_, _, err = b.client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
