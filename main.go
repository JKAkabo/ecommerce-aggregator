package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sqlTest/scrapers"
	"strings"
)

const (
	DB_DSN = "postgresql://johnkennedy:@localhost:5432/rentpal?sslmode=disable"
)

func getAds(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	query := params.Get("query")
	c := make(chan scrapers.SearchResult)
	go scrapers.SearchJiji(query, c)
	go scrapers.SearchJumia(query, c)
	go scrapers.SearchTonaton(query, c)

	var ads []scrapers.Ad
	for i := 0; i < 3; i++ {
		result := <-c
		if result.Err != nil {
			log.Fatal(result.Err)
		}
		ads = append(ads, result.Ads...)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ads)
}

func insertAds(db *sql.DB, ads []scrapers.Ad) error {
	var placeholders []string
	var values []interface{}

	for i := 0; i < len(ads); i++ {
		ad := ads[i]
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", len(values)+1, len(values)+2, len(values)+3, len(values)+4, len(values)+5, len(values)+6))
		values = append(values, ad.Name)
		values = append(values, ad.Price)
		values = append(values, ad.Url)
		values = append(values, ad.Currency)
		values = append(values, ad.Image)
		values = append(values, ad.Platform)
	}
	q := "insert into ads (name, price, url, currency, image, platform) values %s"
	q = fmt.Sprintf(q, strings.Join(placeholders, ","))
	_, err := db.Exec(q, values...)
	return err
}

func main() {
	//db, _ := sql.Open("postgres", DB_DSN)
	router := mux.NewRouter()

	router.HandleFunc("/ads", getAds).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		 port = "8080"
	}
	log.Fatal(http.ListenAndServe(":" + port, router))
}
