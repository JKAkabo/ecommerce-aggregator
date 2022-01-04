package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"math"
	"strconv"
	"strings"
)

type tResponse struct {
	Serp tSerp `json:"serp"`
}

type tSerp struct {
	Ads tAds `json:"ads"`
}

type tAds struct {
	Data tData `json:"data"`
}

type tData struct {
	Ads []tAd `json:"ads"`
	PaginationData tPaginationData `json:"paginationData"`
}

type tPaginationData struct {
	ActivePage int `json:"activePage"`
	Total int `json:"total"`
	PageSize int `json:"pageSize"`
}

type tAd struct {
	Title string `json:"title"`
	ImgUrl string `json:"imgUrl"`
	Price string `json:"price"`
	Slug string `json:"slug"`
}

func parseProductPrice(priceStr string) float64 {
	a := strings.Split(priceStr, " ")
	if len(a) < 2 {
		return 0
	}
	s := a[1]
	s = strings.Replace(s, ",", "", -1)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func getTonatonAds(url string) (ads []Ad, nextUrl string, err error) {
	collector := colly.NewCollector()
	collector.OnHTML("script", func(script *colly.HTMLElement) {
		if strings.Contains(script.Text, "window.initialData = ") {
			j := strings.Split(script.Text, "window.initialData = ")[1]
			d := tResponse{}
			_ = json.Unmarshal([]byte(j), &d)

			for i := 0; i < len(d.Serp.Ads.Data.Ads); i++ {
				ad := Ad{
					Name:     d.Serp.Ads.Data.Ads[i].Title,
					Price:    parseProductPrice(d.Serp.Ads.Data.Ads[i].Price),
					Url:      "https://tonaton.com/en/ad/" + d.Serp.Ads.Data.Ads[i].Slug,
					Currency: "GHS",
					Image:    d.Serp.Ads.Data.Ads[i].ImgUrl,
					Platform: "TONATON",
				}
				ads = append(ads, ad)
			}
			lastPage := int(math.Ceil(float64(d.Serp.Ads.Data.PaginationData.Total) / float64(d.Serp.Ads.Data.PaginationData.PageSize)))
			if d.Serp.Ads.Data.PaginationData.ActivePage < 400 && d.Serp.Ads.Data.PaginationData.ActivePage < lastPage {
				nextUrl = "https://tonaton.com/en/ads/ghana?sort=date&order=desc&buy_now=0&urgent=0&page=" + strconv.Itoa(d.Serp.Ads.Data.PaginationData.ActivePage + 1)
			}
		}
	})
	collector.OnScraped(func(response *colly.Response) {
		log.Println("Scraped", response.Request.URL)
	})

	err = collector.Visit(url)
	collector.Wait()
	if err != nil {
		return nil, "", err
	}
	return ads, nextUrl, nil
}

func SearchTonaton(query string, c chan SearchResult) {
	var ads []Ad
	nextUrl := "https://tonaton.com/en/ads?page=1"
	var err error
	i := 0
	for nextUrl != "" {
		var moreAds []Ad
		moreAds, nextUrl, err = getTonatonAds(nextUrl + "&query=" + query)
		if err != nil {
			c <- SearchResult{err, nil}
			return
		}
		ads = append(ads, moreAds...)

		// page is reset to 1 when the last page is exceeded
		if strings.HasSuffix(nextUrl, "=1") {
			break
		}
		if i++; i > 20 {
			break
		}
	}
	c <- SearchResult{nil, ads}
}

//func SearchTonaton() ([]Ad, error) {
//	var ads []Ad
//	collector := colly.NewCollector()
//	collector.OnRequest(func(request *colly.Request) {
//		log.Println("Visiting", request.URL)
//	})
//	collector.OnResponse(func(response *colly.Response) {
//		if response.StatusCode == 200 {
//			log.Println("Got Response from ", response.Request.URL)
//			response.Save("response.html")
//		}
//	})
//	collector.OnHTML("ul.list--3NxGO", func(containerEl *colly.HTMLElement) {
//		log.Println("Container Found")
//		containerEl.ForEach("li.gtm-normal-ad", func(i int, productEl *colly.HTMLElement) {
//			productName := productEl.ChildText("a > div > div.content--3JNQz > h2")
//			productPrice := parseProductPrice(productEl.ChildText("a > div > div.content--3JNQz > div:first-of-type > div.price--3SnqI > span:first-of-type"))
//			productUrl := "https://tonaton.com/" + productEl.ChildAttr("a", "href")
//			if productName != "" {
//				ad := Ad{Name: productName, Price: productPrice, Url: productUrl, Currency: "GHS"}
//				ads = append(ads, ad)
//			}
//		})
//	})
//	err := collector.Visit("https://tonaton.com/en/ads")
//	collector.Wait()
//	if err != nil {
//		return nil, err
//	}
//	return ads, nil
//}