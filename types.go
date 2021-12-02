package main

import "encoding/xml"

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   items    `xml:"item"`
}

type item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Descr   string   `xml:"description"`
	PubDate string   `xml:"pubDate"`
}

type items []item

func (v items) Len() int {
	return len(v)
}

func (v items) Less(i, j int) bool {
	vi := transformString2Timestamp(v[i].PubDate)
	vj := transformString2Timestamp(v[j].PubDate)
	return vi.Before(vj)
}

func (v items) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
