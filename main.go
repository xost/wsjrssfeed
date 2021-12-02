package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"sort"
)

const sender = "wsjrobot@granfin.ru"
const smtpServer = "mail.granfin.ru:25"

var receiver = []string{"stas@granfin.ru"}

var count = 10.0
var total = 0
var price float64 = 0.0

var sources = map[string]struct {
	link   string
	weight int
	items  *items
}{
	"opinion":     {"https://feeds.a.dj.com/rss/RSSOpinion.xml", 2, nil},
	"worldNews":   {"https://feeds.a.dj.com/rss/RSSWorldNews.xml", 1, nil},
	"usBuiseness": {"https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml", 1, nil},
	"market":      {"https://feeds.a.dj.com/rss/RSSMarketsMain.xml", 1, nil},
	"tech":        {"https://feeds.a.dj.com/rss/RSSWSJD.xml", 1, nil},
	"lifestyle":   {"https://feeds.a.dj.com/rss/RSSLifestyle.xml", 1, nil},
}

var begin = "To: stas@granfin.ru\r\n" +
	"From: <wsjrobot@granfin.ru>\r\n" +
	"Subject: WSJ random articles\r\n"

const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

var msg = "<HTML><BODY>"

var end = "</BODY></HTML>"

func main() {
	for category, url := range sources {
		rawData, err := getRssXml(url.link)
		if err != nil {
			log.Println(err)
			continue
		}
		items := getRssItems(rawData)
		if url.weight > len(items) {
			url.weight = len(items)
		}
		total += url.weight
		for i := range items {
			items[i].PubDate = items[i].PubDate[5:25]
		}
		sort.Sort(items)
		url.items = &items
		sources[category] = url
	}
	price = float64(count) / float64(total)
	for _, url := range sources {
		msg += makeMessage(*url.items, url.weight)
	}
	msg = begin + mime + msg + end
	sendEmail(msg)
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

func getRssItems(rawData []byte) items {
	rss := rss{}
	err := xml.Unmarshal(rawData, &rss)
	if err != nil {
		return nil
	}
	return rss.Channel.Items
}

func makeMessage(items items, weight int) string {
	q := float64(weight) * price
	log.Println(count, weight, q)
	msg := ""
	for i := 0; count > 0 && i < len(items); i++ {
		count -= q
		msg += fmt.Sprintf("<h4><a href='%s'>%s</a></h4>publish date: %s<br/>", items[i].Link, items[i].Title, items[i].PubDate)
	}
	return msg
}

func sendEmail(msg string) {
	smtp.SendMail(
		smtpServer,
		nil,
		sender,
		receiver,
		[]byte([]byte(msg)),
	)
}
