package feishu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	fs, err := NewFeishu("")
	assert.NoError(t, err)
	err = fs.SendMessage([][]MessageBlock{
		{
			{
				Tag:  "text",
				Text: "需要重新登录：",
			},
			{
				Tag:  "a",
				Text: "Link",
				Href: "https://www.baidu.com",
			},
		},
	})
	assert.NoError(t, err)
}
