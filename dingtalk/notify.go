package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DingTalk struct {
	webhook    string
	httpClient *http.Client
}

func NewDingTalk(accessToken string) (*DingTalk, error) {
	webhook := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken)
	return &DingTalk{
		webhook:    webhook,
		httpClient: http.DefaultClient,
	}, nil
}

func (self *DingTalk) SendMessage(msg *Message) error {
	data, _ := json.Marshal(msg)
	req, err := http.NewRequest("POST", self.webhook, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := self.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doHTTPRequest: %v", err)
	}

	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("StatusCode: %d, %s", resp.StatusCode, string(content))
	}

	r := new(Response)
	if err = json.Unmarshal(content, r); err != nil {
		return fmt.Errorf("read body: %s: %v", string(content), err)
	}
	if r.ErrCode != 0 {
		return fmt.Errorf("errcode: %d, errmsg: %s", r.ErrCode, r.ErrMsg)
	}

	return nil
}
