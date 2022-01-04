package scrapers

type Ad struct {
	Name string `json:"name"`
	Price float64 `json:"price"`
	Url string `json:"url"`
	Currency string `json:"currency"`
	Image string `json:"image"`
	Platform string `json:"platform"`
}

type SearchResult struct {
	Err error
	Ads []Ad
}