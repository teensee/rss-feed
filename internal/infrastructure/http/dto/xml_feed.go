package dto

import (
	"encoding/xml"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	XMLNSDC string   `xml:"xmlns:dc,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	PubDate     string  `xml:"pubDate"`
	Image       *Image  `xml:"image"`
	Items       *[]Item `xml:"item"`
}

type Image struct {
	Link  string `xml:"link"`
	URL   string `xml:"url"`
	Title string `xml:"title"`
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Creator     string   `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Categories  []string `xml:"category"`
}

type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}
