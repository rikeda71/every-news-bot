package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/s14t284/every-news-bot/service"
)

// FlexMessageに利用するJSONテンプレート
const (
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
				"label": "ニュースを見る",
				"uri": "%s"
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
	url, err := DecideRequestURL(message)
	if err != nil {
		return url, err
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

// GetLineBotHandler : Linebotの設定を行うメソッド
func GetLineBotHandler() (*httphandler.WebhookHandler, error) {
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
