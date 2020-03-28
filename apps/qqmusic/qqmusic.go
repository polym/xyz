package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/polym/xyz/ffmpeg"
	"github.com/upyun/go-sdk/upyun"
)

var (
	urlChan      chan string
	musicPath    string
	waitDuration = time.Duration(30) * time.Second
	up           *upyun.UpYun
)

func ServeHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = r.Host
			urlChan <- req.URL.String()
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}

func doJob() {
	for url := range urlChan {
		log.Printf("url %s", url)
		mediaInfo, err := ffmpeg.Probe(url, "", waitDuration)
		if err != nil {
			log.Printf("ffprobe %s: %v", url, err)
			continue
		}

		tags := mediaInfo.Format.Tags
		if tags.Album == "" {
			tags.Album = "unknown"
		}
		filename := filepath.Join(musicPath, tags.Artist, tags.Album, tags.Title+".mp3")
		if _, err = os.Stat(filename); !os.IsNotExist(err) {
			if fInfo, _ := up.GetInfo(filename); fInfo != nil {
				log.Printf("%s already exists", filename)
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			log.Printf("mkdir %s: %v", filename, err)
			continue
		}

		err = ffmpeg.FFmpeg(ffmpeg.FFmpegOption{
			InputMedia:  url,
			OutputMedia: filename,
		}, waitDuration)
		if err != nil {
			log.Printf("ffmpeg error %s: %v", filename, err)
			continue
		}

		log.Printf("download %s OK", filename)
		err = up.Put(&upyun.PutObjectConfig{
			Path:      filename,
			LocalPath: filename,
		})
		log.Printf("upload %s to %s %v", filename, filename, err)
	}
}

func main() {
	urlChan = make(chan string, 1024)
	musicPath = "music"
	up = upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "icloud01",
		Operator: "myworker",
		Password: "TYGHBNtyghbn",
	})
	go doJob()
	ServeHTTP()
}
