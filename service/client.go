package service

import (
	"encoding/xml"
	"net/http"

	"github.com/s14t284/every-news-bot/entities"
)

// RssReader interface
type RssReader interface {
	Request(url string) (*entities.Channel, error)
}

// RssClient definition
type RssClient struct {
	URL string
}

// Request Yahoo News Rss
func (rc *RssClient) Request(url string) (*entities.Channel, error) {
	rssResult := entities.Channel{}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := xml.NewDecoder(resp.Body).Decode(&rssResult); err != nil {
		return nil, err
	}
	return &rssResult, nil
}
