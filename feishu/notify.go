package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MessageBlock struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
	Href string `json:"href"`
}

type Message [][]MessageBlock

type Feishu struct {
	webhook    string
	httpClient *http.Client
}

func NewFeishu(webhook string) (*Feishu, error) {
	return &Feishu{
		webhook:    webhook,
		httpClient: http.DefaultClient,
	}, nil
}

func (fs *Feishu) SendTextMessage(text string) error {
	return fs.SendMessage([][]MessageBlock{
		{
			{
				Tag:  "text",
				Text: text,
			},
		},
	})
}

func (fs *Feishu) SendMessage(msg Message) error {
	data, _ := json.Marshal(map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"content": msg,
				},
			},
		},
	})
	req, err := http.NewRequest("POST", fs.webhook, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := fs.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doHTTPRequest: %v", err)
	}

	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(content))
	}

	return nil
}
