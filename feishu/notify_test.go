package feishu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	fs, err := NewFeishu("")
	assert.NoError(t, err)
	err = fs.SendMessage("Hello")
	assert.NoError(t, err)
}
