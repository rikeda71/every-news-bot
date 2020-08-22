package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/s14t284/news-bot-lambda/service"
)

// 主要カテゴリのニュース
const (
	Main                     = "https://news.yahoo.co.jp/pickup/rss.xml"               // 主要
	Domestic                 = "https://news.yahoo.co.jp/pickup/domestic/rss.xml"      // 国内
	World                    = "https://news.yahoo.co.jp/pickup/world/rss.xml"         // 国際
	Economy                  = "https//news.yahoo.co.jp/pickup/economy/rss.xml"        // 経済
	Entertainment            = "https://news.yahoo.co.jp/pickup/entertainment/rss.xml" // エンタメ
	Sports                   = "https://news.yahoo.co.jp/pickup/sports/rss.xml"        // スポーツ
	Computer                 = "https://news.yahoo.co.jp/pickup/computer/rss.xml"      // IT
	Science                  = "https://news.yahoo.co.jp/pickup/science/rss.xml"       // 科学
	Local                    = "https://news.yahoo.co.jp/pickup/local/rss.xml"         // 地域
	LineFlexMessageObjFormat = `{
		"type": "bubble",
		"header": {
		  "type": "box",
		  "layout": "baseline",
		  "contents": [
			{
			  "type": "text",
			  "text": "%s",
			  "size": "lg",
			  "color": "#333333",
			  "weight": "bold",
			  "position": "absolute"
			}
		  ]
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "button",
			  "action": {
				"type": "uri",
				"label": "%s",
				"uri": "http://linecorp.com/"
			  }
			}
		  ]
		}
	  }`
	LineMessageJSONFormat = `{"type": "carousel", "contents":[%s]}`
)

func getResponseObject(url string) ([]string, error) {
	var client service.RssReader = &service.RssClient{}
	var newsObjects []string
	resp, err := client.Request(url)
	if err != nil {
		return newsObjects, err
	}
	// LineBotのFlexMessageの形式で格納していく
	for _, news := range resp.Items {
		s := fmt.Sprintf(LineFlexMessageObjFormat, news.Title, news.Link)
		newsObjects = append(newsObjects, s)
	}
	return newsObjects, nil
}

func getReplyTextWithNewsCategory(message string) (string, error) {
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
		return "不正な文字列です", errors.New("不正な文字列です")
	}
	newsObj, err := getResponseObject(url)
	if err == nil {
		return fmt.Sprintf(LineMessageJSONFormat, strings.Join(newsObj, ",")), nil
	}
	return "リクエストに失敗しました", err
}

func replyMessageSetting(event *linebot.Event, bot *linebot.Client) {
	replyWhenInvalid := func(reply string, err error) {
		log.Print(err)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			log.Print(err)
		}

	}
	if event.Type == linebot.EventTypeMessage {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			reply, err := getReplyTextWithNewsCategory(message.Text)
			if err != nil {
				// ニュースカテゴリ以外の文字列の場合
				replyWhenInvalid(reply, err)
			} else {
				// ニュースカテゴリの場合は、FlexMessageで表示
				container, err := linebot.UnmarshalFlexMessageJSON([]byte(reply))
				if err != nil {
					replyWhenInvalid(reply, err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("alt text", container)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

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
			replyMessageSetting(event, bot)
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
