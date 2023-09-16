package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFedd struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (*RSSFedd, error) {
	httpClinet := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClinet.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := &RSSFedd{}
	err = xml.Unmarshal(data, rssFeed)
	if err != nil {
		return nil, err
	}

	return rssFeed, err
}
