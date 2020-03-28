package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/polym/xyz/ffmpeg"
)

func main() {
	var (
		filename   string
		httpPrefix string
	)

	flag.StringVar(&filename, "f", "", "ts filename")
	flag.StringVar(&httpPrefix, "p", "", "http url prefix")
	flag.Parse()

	fd, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open %s: %v", filename, err)
		return
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	durations := []float64{}
	TARGETDURATION := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http") {
			line = httpPrefix + line
		}
		mediaInfo, err := ffmpeg.Probe(line, "", time.Duration(60)*time.Second)
		if err != nil {
			fmt.Printf("ffprobe %s: %v\n", line, err)
			return
		}

		dur, err := strconv.ParseFloat(mediaInfo.Format.Duration, 64)
		if err != nil {
			fmt.Printf("invalid duration\n")
			return
		}

		lines = append(lines, line)
		durations = append(durations, dur)
		if int(dur+0.5) > TARGETDURATION {
			TARGETDURATION = int(dur + 0.5)
		}
	}

	fmt.Printf("#EXTM3U\n")
	fmt.Printf("#EXT-X-VERSION:3\n")
	fmt.Printf("#EXT-X-MEDIA-SEQUENCE:1\n")
	fmt.Printf("#EXT-X-TARGETDURATION:%d\n", TARGETDURATION)
	for idx, line := range lines {
		fmt.Printf("#EXTINF:%.3f,\n", durations[idx])
		fmt.Printf("%s\n", path.Base(line))
	}
	fmt.Printf("#EXT-X-ENDLIST\n")
}
