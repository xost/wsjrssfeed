package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []item   `xml:"item"`
}

type item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Descr   string   `xml:"description"`
	PubDate string   `xml:"pubDate"`
}

var sources = map[string]string{
	"opinion":     "https://feeds.a.dj.com/rss/RSSOpinion.xml",
	"worldNews":   "https://feeds.a.dj.com/rss/RSSWorldNews.xml",
	"usBuiseness": "https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml",
	"market":      "https://feeds.a.dj.com/rss/RSSMarketsMain.xml",
	"tech":        "https://feeds.a.dj.com/rss/RSSWSJD.xml",
	"lifestyle":   "https://feeds.a.dj.com/rss/RSSLifestyle.xml",
}

func main() {
	for _, url := range sources {
		rawData, err := getRssXml(url)
		if err != nil {
			log.Println(err)
			continue
		}
		items := getRssItems(rawData)
		for _, i := range items {
			log.Println(i.Title)
		}
	}
}

func getRssXml(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rawData, nil
}

func getRssItems(rawData []byte) []item {
	rss := rss{}
	err := xml.Unmarshal(rawData, &rss)
	if err != nil {
		return nil
	}
	return rss.Channel.Items
}

func makeMessage(items []item) {
}
