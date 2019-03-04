package main

import (
	"fmt"
	"go-slack-bot/myslack"
	"os"

	"go-slack-bot/config"
	"go-slack-bot/queue"
	"log"

	"github.com/go-chat-bot/bot"

	"github.com/go-chat-bot/bot/slack"
)

var token string

func init() {
	token = os.Getenv("SLACK_TOKEN")
	if token == "" {
		token = "xoxb-9329002610-389158948503-Rjx15Rh2b2aNDauY6oap0IzK"
	}

	bot.RegisterCommand("hello",
		"Sends a 'Hello' message to you on the channel.",
		"",
		hello)

	bot.RegisterCommand("queue",
		"Get RabbitMQ list Queues.",
		"",
		queueCmd)
}

func hello(command *bot.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("Hello %s", command.User.RealName)
	return
}

func queueCmd(command *bot.Cmd) (msg string, err error) {
	instance := &myslack.MySlack{
		Token: token,
	}
	err = queue.GetQueues(command.Channel, instance)
	return "", err
}

func main() {
	slack.Run(token)
	log.Println(config.Setting.App.Queue.Host)
}
