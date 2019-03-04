package myslack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	StyleSuccess = "#5cb85c"
	StyleInfo    = "#5bc0de"
	StyleWarning = "#f0ad4e"
	StyleError   = "#d9534f"
)

// MySlack is instance
type MySlack struct {
	Token string
}

// ResponseMessage struct for message option
type ResponseMessage struct {
	Type        string        `json:"type"`
	Channel     string        `json:"channel"`
	Attachments []Attachments `json:"attachments"`
	Username    string        `json:"username"`
}

// Attachments struct for message option
type Attachments struct {
	Color    string   `json:"color"`
	Title    string   `json:"title"`
	Pretext  string   `json:"pretext"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

// Message is
func (s *MySlack) Message(response *ResponseMessage) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	req, err := http.NewRequest(http.MethodPost, "https://api.slack.com/api/chat.postMessage", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	client := &http.Client{}
	// add header basic authen
	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http %d , body : %s", resp.StatusCode, string(body))
	}

	fmt.Println(string(body))

	return nil
}
