package rss

type Feed struct {
	title       string
	link        string
	description string
	pubDate     string
	image       *Image
	items       []*Item
}

func (f *Feed) GetTitle() string {
	return f.title
}

func (f *Feed) GetLink() string {
	return f.link
}

func (f *Feed) GetDescription() string {
	return f.description
}

func (f *Feed) GetPubDate() string {
	return f.pubDate
}

func (f *Feed) GetImage() *Image {
	return f.image
}

func (f *Feed) GetItems() []*Item {
	return f.items
}

func NewFeed(
	title string,
	link string,
	description string,
	pubDate string,
	image *Image,
	items []*Item,
) *Feed {
	return &Feed{
		title:       title,
		link:        link,
		description: description,
		pubDate:     pubDate,
		image:       image,
		items:       items,
	}
}

type Image struct {
	link  string
	url   string
	title string
}

func NewImage(link, url, title string) *Image {
	return &Image{
		link:  link,
		url:   url,
		title: title,
	}
}

func (i *Image) GetLink() string {
	return i.link
}

func (i *Image) GetUrl() string {
	return i.url
}

func (i *Image) GetTitle() string {
	return i.title
}

type Item struct {
	title       string
	link        string
	description string
	pubDate     string
	creator     string
	categories  []string
}

func NewItem(
	title string,
	link string,
	description string,
	pubDate string,
	creator string,
	categories []string,
) *Item {
	return &Item{
		title:       title,
		link:        link,
		description: description,
		pubDate:     pubDate,
		creator:     creator,
		categories:  categories,
	}
}

func (i *Item) GetTitle() string {
	return i.title
}

func (i *Item) GetLink() string {
	return i.link
}
func (i *Item) GetDescription() string {
	return i.description
}

func (i *Item) GetPubDate() string {
	return i.pubDate
}

func (i *Item) ChangePubDate(pubDate string) {
	i.pubDate = pubDate
}

func (i *Item) GetCreator() string {
	return i.creator
}

func (i *Item) GetCategories() []string {
	return i.categories
}
