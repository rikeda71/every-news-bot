package utils

import "errors"

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

// DecideRequestURL : メッセージからリクエスト先のURLを決定し返却する
func DecideRequestURL(message string) (string, error) {
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
	return url, nil
}
