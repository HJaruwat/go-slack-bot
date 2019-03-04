package queue

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-slack-bot/config"
	"go-slack-bot/myslack"
	"io/ioutil"
	"net/http"
)

var basic string
var setting *config.Queue

type (

	//queueDetail struct of
	queueDetail struct {
		Name         string `json:"name"`
		Message      int    `json:"messages"`
		MessageReady int    `json:"messages_ready"`
		Consumers    int    `json:"consumers"`
	}

	Message struct {
		Attachment []Attachment `json:"attachments"`
	}

	Attachment struct {
		Text  string `json:"text"`
		Color string `json:"color"`
	}
)

func init() {
	setting = &config.Setting.App.Queue
	str := fmt.Sprintf("%s:%s", setting.Username, setting.Password)
	basic = base64.StdEncoding.EncodeToString([]byte(str))
}

// GetQueues is
func GetQueues(channel string, slack *myslack.MySlack) error {
	attachments := []myslack.Attachments{}
	response := &myslack.ResponseMessage{
		Type:        "message",
		Channel:     channel,
		Attachments: attachments,
		Username:    "bot-backend",
	}

	attachment := myslack.Attachments{
		Color:    myslack.StyleSuccess,
		Pretext:  "Check message in rabbitMQ",
		Text:     "Everything is *ok!*",
		Title:    "RabbitMQ",
		MrkdwnIn: []string{"text", "pretext"},
	}

	endpoint := "http://" + setting.Tool + "/api/queues/%2F"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	// add header basic authen
	req.Header.Add("Authorization", "Basic "+basic)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Http request error (%d) : %s", resp.StatusCode, string(body))
	}

	details := []queueDetail{}
	err = json.Unmarshal(body, &details)
	if err != nil {
		return err
	}

	if len(details) > 0 {
		var out string
		var ok bool
		var counter = 1
		ok = true

		for _, d := range details {
			if d.Message > 0 || d.MessageReady > 0 {
				ok = false
				color := getColor(d.Message, d.Consumers)
				text := fmt.Sprintf("%d. %s `%-30s` message : `%d` worker : `%d`\n", counter, color, d.Name, d.Message, d.Consumers)
				out += text
				counter++
			}
		}

		if !ok {
			attachment.Color = myslack.StyleWarning
			attachment.Text = out
		}

		attachments = append(attachments, attachment)
		response.Attachments = attachments
		return slack.Message(response)

	}

	attachment.Color = myslack.StyleWarning
	attachment.Text = "*Queue not found*"
	attachments = append(attachments, attachment)

	response.Attachments = attachments
	return slack.Message(response)
}

//getColor is function
func getColor(message int, worker int) string {
	if message > 0 && worker < 1 {
		return ":rage:"
	} else if message > 0 && worker > 0 {
		return ":fearful:"
	} else {
		return ":trollface:"
	}
}
