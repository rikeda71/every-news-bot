package main

import (
	"fmt"

	"github.com/s14t284/news-bot-lambda/service"
)

// 主要カテゴリのニュース
const (
	RssURL = "https://news.yahoo.co.jp/pickup/rss.xml"
)

func main() {
	var client service.RssReader = &service.RssClient{}
	if resp, err := client.Request(RssURL); err == nil {
		fmt.Println(resp)
	}
}
