package entities

import "encoding/xml"

// Channel : Rss全体を表す構造体
type Channel struct {
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	Language    string `xml:"channel>language"`
	PubDate     string `xml:"channel>pubDate"`
	Items       []Item `xml:"channel>item"`
}

// Item : 1ニュースを表す構造体
type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	PubDate string   `xml:"pubDate"`
	GUID    string   `xml:"guid"`
}
