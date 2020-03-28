package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	"github.com/polym/xyz/dingtalk"
)

func main() {
	var token string
	flag.StringVar(&token, "tk", "", "dingding access token")
	flag.Parse()
	dtalk, _ := dingtalk.NewDingTalk(token)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		err := dtalk.SendMessage(dingtalk.NewTextMessage(line))
		log.Printf("send message: %v\n", err)
	}

	for {
		time.Sleep(time.Minute)
	}
}
