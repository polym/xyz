package ffmpeg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/polym/xyz/exec"
)

type StreamInfo struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	BitRate   string `json:"bit_rate"`

	// video
	VideoWidth              int    `json:"width,omitempty"`
	VideoHeight             int    `json:"height,omitempty"`
	VideoDisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`

	// audio
	AudioSampleFmt  string `json:"sample_fmt,omitempty"`
	AudioSampleRate string `json:"sample_rate,omitempty"`
}

type FormatTags struct {
	Title  string `json:"title,omitempty"`
	Album  string `json:"album,omitempty"`
	Artist string `json:"artist,omitempty"`
}

type FormatInfo struct {
	Filename string      `json:"filename"`
	Format   string      `json:"format_name"`
	Duration string      `json:"duration"`
	Tags     *FormatTags `json:"tags,omitempty"`
}

type MediaInfo struct {
	Streams []*StreamInfo `json:"streams"`
	Format  *FormatInfo   `json:"format"`
}

func (m *MediaInfo) String() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

func Probe(media string, headers string, timeout time.Duration) (*MediaInfo, error) {
	args := []string{
		"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams",
	}
	if headers != "" {
		args = append(args, "-headers", headers)
	}
	args = append(args, media)
	stdout, stderr, err := exec.Command("ffprobe", args...).GetOutputWithTimeout(timeout)
	if err != nil {
		return nil, fmt.Errorf("ffprobe %s: %v %s", media, err, string(stderr))
	}

	mediaInfo := &MediaInfo{}
	err = json.Unmarshal(stdout, mediaInfo)
	if err != nil {
		return nil, fmt.Errorf("json decode media info: %v", err)
	}
	return mediaInfo, nil
}
