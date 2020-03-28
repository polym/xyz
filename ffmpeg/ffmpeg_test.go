package ffmpeg

import (
	"fmt"
	"testing"
)

func TestFFmpeg(t *testing.T) {
	opt := FFmpegOption{
		InputMedia:  Media,
		OutputMedia: "1.mkv",
	}
	fmt.Println(FFmpeg(opt, -1))
}
