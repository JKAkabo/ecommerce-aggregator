package scrapers

import (
	"github.com/mxschmitt/playwright-go"
	"strconv"
	"strings"
)

func SearchJumia(query string, c chan SearchResult) {
	pw, err := playwright.Run()
	if err != nil {
		c <- SearchResult{err, nil}
		return
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		c <- SearchResult{err, nil}
		return
	}
	page, err := browser.NewPage()
	page.SetExtraHTTPHeaders(map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36"})
	if err != nil {
		c <- SearchResult{err, nil}
		return
	}
	if _, err := page.Goto("https://www.jumia.com.gh/catalog/?viewType=grid&q=" + query + "#catalog-listing"); err != nil {
		c <- SearchResult{err, nil}
		return
	}

	products, _ := page.QuerySelectorAll("article.c-prd")
	var ads []Ad
	for _, product := range products {
		anchorEl, _ := product.QuerySelector("a")
		href, _ := anchorEl.GetAttribute("href")
		if href == "" {
			continue
		}
		nameEl, _ := product.QuerySelector("a > div.info > h3.name")
		name, _ := nameEl.InnerText()
		priceEl, _ := product.QuerySelector("a > div.info > div.prc")
		priceStr, _ := priceEl.InnerText()
		price, _ := strconv.ParseFloat(strings.Split(priceStr, " ")[1], 64)

		imageEl, _ := product.QuerySelector("a > div.img-c > img")
		image, _ := imageEl.GetAttribute("src")
		ad := Ad{
			Name:     name,
			Price:    price,
			Url:      "https://www.jumia.com.gh/" + href,
			Currency: "GHS",
			Image:    image,
			Platform: "JUMIA",
		}
		ads = append(ads, ad)
	}
	c <- SearchResult{nil, ads}
}
