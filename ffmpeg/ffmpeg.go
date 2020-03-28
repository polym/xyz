package ffmpeg

import (
	"fmt"
	"time"

	"github.com/polym/xyz/exec"
)

type VideoOption struct {
	Codec   string
	Scale   string
	BitRate int
	Disable bool
}

func (v *VideoOption) ToArgs() (args []string) {
	if v.Disable {
		args = append(args, "-vn")
		return
	}
	if v.Codec != "" {
		args = append(args, "-vcodec", v.Codec)
	}
	if v.Scale != "" {
		args = append(args, "-vf", fmt.Sprintf("scale=%s", v.Scale))
	}
	if v.BitRate != 0 {
		args = append(args, "-vb", fmt.Sprint(v.BitRate))
	}
	return
}

type AudioOption struct {
	Codec   string
	Disable bool
}

func (a *AudioOption) ToArgs() (args []string) {
	if a.Disable {
		args = append(args, "-an")
		return
	}
	if a.Codec != "" {
		args = append(args, "-acodec", a.Codec)
	}
	return
}

type FFmpegOption struct {
	InputMedia string

	VideoOpt *VideoOption
	AudioOpt *AudioOption

	OutputFormat string
	OutputMedia  string

	LogLevel string

	AdvancedArgs []string
}

func (f FFmpegOption) ToArgs() (args []string) {
	if f.LogLevel == "" {
		args = append(args, "-v", "quiet")
	} else {
		args = append(args, "-v", f.LogLevel)
	}

	args = append(args, "-i", f.InputMedia)

	if f.VideoOpt != nil {
		args = append(args, f.VideoOpt.ToArgs()...)
	}

	if f.AudioOpt != nil {
		args = append(args, f.AudioOpt.ToArgs()...)
	}

	if len(f.AdvancedArgs) > 0 {
		args = append(args, f.AdvancedArgs...)
	}

	if f.OutputFormat != "" {
		args = append(args, "-f", f.OutputFormat)
	}

	if f.OutputMedia != "" {
		args = append(args, f.OutputMedia)
	}

	return
}

func FFmpeg(opt FFmpegOption, timeout time.Duration) error {
	args := opt.ToArgs()
	_, stderr, err := exec.Command("ffmpeg", args...).GetOutputWithTimeout(timeout)
	if err != nil {
		return fmt.Errorf("ffmpeg %v: %s", err, string(stderr))
	}
	return nil
}
