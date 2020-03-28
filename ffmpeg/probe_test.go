package ffmpeg

import (
	"fmt"
	"testing"
)

var (
	Media = "http://bigfile.b0.upaiyun.com/x.mp4"
)

func TestProbe(t *testing.T) {
	_, err := Probe(Media, "", -1)
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
}
