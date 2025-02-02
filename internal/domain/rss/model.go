package rss

type Feed struct {
	XMLName string
	Version string
	XMLNSDC string
	Channel *Channel
}

type Channel struct {
	Title       string
	Link        string
	Description string
	PubDate     string
	Image       *Image
	Items       []*Item
}

type Image struct {
	link  string
	uRL   string
	title string
}

func NewImage(link string, URL string, title string) *Image {
	return &Image{
		link:  link,
		uRL:   URL,
		title: title,
	}
}

type Item struct {
	Title       string
	Link        string
	Description string
	PubDate     string
	Creator     string
	Categories  []string
}
