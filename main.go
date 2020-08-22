package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/s14t284/news-bot-lambda/service"
)

// 主要カテゴリのニュース
const (
	Main          = "https://news.yahoo.co.jp/pickup/rss.xml"               // 主要
	Domestic      = "https://news.yahoo.co.jp/pickup/domestic/rss.xml"      // 国内
	World         = "https://news.yahoo.co.jp/pickup/world/rss.xml"         // 国際
	Economy       = "https//news.yahoo.co.jp/pickup/economy/rss.xml"        // 経済
	Entertainment = "https://news.yahoo.co.jp/pickup/entertainment/rss.xml" // エンタメ
	Sports        = "https://news.yahoo.co.jp/pickup/sports/rss.xml"        // スポーツ
	Computer      = "https://news.yahoo.co.jp/pickup/computer/rss.xml"      // IT
	Science       = "https://news.yahoo.co.jp/pickup/science/rss.xml"       // 科学
	Local         = "https://news.yahoo.co.jp/pickup/local/rss.xml"         // 地域
)

func getResponseText(url string) (string, error) {
	var client service.RssReader = &service.RssClient{}
	text := ""
	resp, err := client.Request(url)
	if err != nil {
		return "不正な文字列です。", err
	}
	for _, news := range resp.Items {
		text += news.Title + "\n"
		text += news.Link + "\n"
	}
	return text, nil
}

func getReplyTextWithNewsCategory(message string) string {
	var url string
	switch message {
	case "主要":
		url = Main
	case "国内":
		url = Domestic
	case "国際":
		url = World
	case "経済":
		url = Economy
	case "エンタメ":
		url = Entertainment
	case "スポーツ":
		url = Sports
	case "IT":
		url = Computer
	case "科学":
		url = Science
	case "地域":
		url = Local
	default:
		url = ""
	}
	log.Print(message)
	log.Print(url)
	text, _ := getResponseText(url)
	return text
}

func getLineHandler() (*httphandler.WebhookHandler, error) {
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					text := getReplyTextWithNewsCategory(message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	return handler, nil
}

func main() {
	handler, err := getLineHandler()
	if err != nil {
		return
	}
	http.Handle("/callback", handler)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
