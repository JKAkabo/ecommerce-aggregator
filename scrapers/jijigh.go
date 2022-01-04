package scrapers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type jResponse struct {
	AdultOnly   bool         `json:"adult_only"`
	AdvertsList jAdvertsList `json:"adverts_list"`
	NextUrl     string       `json:"next_url"`
}

type jAdvertsList struct {
	Adverts []jAdvert `json:"adverts"`
}

type jAdvert struct {
	Title    string    `json:"title"`
	PriceObj jPriceObj `json:"price_obj"`
	ImageObj jImageObj `json:"image_obj"`
	Url      string    `json:"url"`
}

type jPriceObj struct {
	Value float64 `json:"value"`
}

type jImageObj struct {
	Url string `json:"url"`
}

func getJijiAds(url string) (ads []Ad, nextUrl string, err error) {
	client := http.Client{}
	res, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	resBody := jResponse{}
	err = json.Unmarshal(bytes, &resBody)
	if err != nil {
		return nil, "", err
	}

	for i := 0; i < len(resBody.AdvertsList.Adverts); i++ {
		advert := resBody.AdvertsList.Adverts[i]
		ad := Ad{
			Name: advert.Title,
			Price: advert.PriceObj.Value,
			Url: advert.Url,
			Currency: "GHS",
			Image: advert.ImageObj.Url,
			Platform: "JIJI",
		}
		ads = append(ads, ad)
	}
	log.Println("Scraped", url)
	return ads, resBody.NextUrl, nil
}

func SearchJiji(query string, c chan SearchResult) {
	var ads []Ad
	// use trailing space to prevent api from redirecting
	nextUrl := "https://jiji.com.gh/api_web/v1/listing?page=1&query=" + query + " "
	var err error
	i := 0
	for nextUrl != "" {
		var moreAds []Ad
		moreAds, nextUrl, err = getJijiAds(nextUrl)
		if err != nil {
			c <- SearchResult{err, nil}
			return
		}
		ads = append(ads, moreAds...)
		if i++; i > 20 {
			break
		}
	}
	c <- SearchResult{nil, ads}
}
