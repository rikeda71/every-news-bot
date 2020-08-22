package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
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
	// var client service.RssReader = &service.RssClient{}
	// if resp, err := client.Request(Main); err == nil {
	// 	for _, news := range resp.Items {
	// 		fmt.Println(news.Title)
	// 		fmt.Println(news.Link)
	// 	}
	// }
}
