package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func (fs *Feishu) SendMessage(text string) error {
	data, _ := json.Marshal(map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
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
