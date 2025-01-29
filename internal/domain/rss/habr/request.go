package habr

import (
	"fmt"
	"strings"
)

// compile time checks
var _ RssFilter = &FeedUrl{}
var _ RssFilter = &BestFeedUrl{}

type RssFilter interface {
	BuildRssUrl() string
}

type FeedUrl struct {
	Rating EnumRating
	Level  EnumLevel
}

func NewNewFeedUrl(rating EnumRating, level EnumLevel) RssFilter {
	return &FeedUrl{Rating: rating, Level: level}
}

func (f *FeedUrl) BuildRssUrl() string {
	b := strings.Builder{}
	b.Write([]byte("rss/articles"))

	if f.Rating > AllRated {
		b.Write([]byte(fmt.Sprintf("/%s", f.Rating.String())))
	}

	if f.Level > AllLevel {
		b.Write([]byte(fmt.Sprintf("/%s", f.Level.String())))
	}

	return strings.ToLower(b.String())
}

type BestFeedUrl struct {
	Period EnumPeriod
	Level  EnumLevel
}

func NewBestFeedUrl(period EnumPeriod, level EnumLevel) RssFilter {
	return &BestFeedUrl{Period: period, Level: level}
}

func (f *BestFeedUrl) BuildRssUrl() string {
	b := strings.Builder{}
	b.Write([]byte("rss/articles/top"))

	b.Write([]byte(fmt.Sprintf("/%s", f.Period.String())))

	if f.Level > AllLevel {
		b.Write([]byte(fmt.Sprintf("/%s", f.Level.String())))
	}

	return strings.ToLower(b.String())
}

//
//type HubFeedUrl struct {
//	BaseFeedUrl
//	Hub    string
//	Rating EnumRating
//}
//
//func NewHubFeedUrl(hub string, rating EnumRating, level EnumLevel) *HubFeedUrl {
//	return &HubFeedUrl{
//		BaseFeedUrl: BaseFeedUrl{Level: level},
//		Rating:      rating,
//		Hub:         hub,
//	}
//}
//
//func (f *HubFeedUrl) BuildRssUrl() string {
//	b := f.buildBaseUrl(fmt.Sprintf("rss/hub/%s", f.Hub))
//
//	if f.Rating > AllRated {
//		b.WriteString(fmt.Sprintf("/%s", f.Rating.String()))
//	}
//
//	return b.String()
//}
//
//type HubBestFeedUrl struct {
//	BestFeedUrl
//	Hub    string
//	Rating EnumRating
//}
//
//func NewHubBestFeedUrl(hub string, period EnumPeriod, level EnumLevel) *HubBestFeedUrl {
//	return &HubBestFeedUrl{
//		BestFeedUrl: BestFeedUrl{BaseFeedUrl{Level: level}, period},
//		Hub:         hub,
//	}
//}
//
//func (f *HubBestFeedUrl) BuildRssUrl() string {
//	b := f.buildBaseUrl(fmt.Sprintf("rss/hub/%s", f.Hub))
//
//	b.WriteString(fmt.Sprintf("/%s", f.Period.String()))
//
//	return b.String()
//}

//// Enum Level

type EnumLevel int

const (
	AllLevel EnumLevel = iota
	Easy
	Normal
	Hard
)

func (l EnumLevel) String() string {
	return [...]string{"All", "Easy", "Normal", "Hard"}[l]
}

//// Enum Rating

type EnumRating int

const (
	AllRated EnumRating = iota
	Rate0
	Rate10
	Rate25
	Rate50
	Rate100
)

func (l EnumRating) String() string {
	return [...]string{"All", "rated0", "rated10", "rated25", "rated50", "rated100"}[l]
}

//// Enum Period

type EnumPeriod int

const (
	AllPeriod EnumPeriod = iota
	Day
	Week
	Month
	Year
)

func (l EnumPeriod) String() string {
	return [...]string{"alltime", "daily", "weekly", "monthly", "yearly"}[l]
}
