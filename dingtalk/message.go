package dingtalk

type MessageType string

const (
	MessageLink     MessageType = "link"
	MessageText     MessageType = "text"
	MessageMarkdown MessageType = "markdown"
)

type Text struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Link struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type Message struct {
	MsgType  MessageType `json:"msgtype"`
	Text     *Text       `json:"text"`
	Markdown *Markdown   `json:"markdown"`
	Link     *Link       `json:"link"`
}

type Response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func (msg Message) isValid() bool {
	if (msg.MsgType == MessageText && msg.Text != nil) ||
		(msg.MsgType == MessageMarkdown && msg.Markdown != nil) ||
		(msg.MsgType == MessageLink && msg.Link != nil) {
		return true
	}
	return false
}

func NewTextMessage(text string) *Message {
	return &Message{
		MsgType:  MessageText,
		Text:     &Text{Content: text},
		Markdown: nil,
		Link:     nil,
	}
}
