package dto

type RssFeedItemProcess struct {
	Rss     string   `json:"rss"`
	Filters []string `json:"filters"`
}

type AppRssFeedRequest struct {
	Items []*RssFeedItemProcess `json:"items"`
}
