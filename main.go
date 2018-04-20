package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Specs    []string `json:"specs"`
	Keyboard string   `json:"keyboard"`
	Price    string   `json:"price"`
	URL      string   `json:"url"`
}

func main() {
	doc, err := goquery.NewDocument("https://www.apple.com/jp/shop/browse/home/specialdeals/mac/macbook_pro/13")
	if err != nil {
		panic(err)
	}

	reJISKeyboard := regexp.MustCompile(`JIS.*キーボード`)
	reUSKeyboard := regexp.MustCompile(`US.*キーボード`)

	var products []*Product
	doc.Find("tr.product").Each(func(i int, s *goquery.Selection) {
		specsSelection := s.Find("td.specs")

		productID := specsSelection.Find("h3 > a").AttrOr("data-relatedlink", "UNKNOWN")
		productID = productID[4:(4 + 8)]
		productID = strings.Replace(productID, "_", "/", -1)

		productName := specsSelection.Find("h3").Text()
		productName = strings.TrimSpace(productName)

		spec := specsSelection.Text()
		spec = strings.Replace(spec, productName, "", -1)

		var specs []string
		for _, d := range strings.Split(spec, "\n") {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}

			specs = append(specs, d)
		}

		price := s.Find("span[itemprop=price]").Text()
		price = strings.TrimSpace(price)

		keyboard := "UNKNOWN"
		url := specsSelection.Find("h3 > a").AttrOr("href", "#")
		if url != "#" {
			url = "https://www.apple.com" + url

			detailDoc, _ := goquery.NewDocument(url)
			detailFullText := detailDoc.Text()
			if reJISKeyboard.MatchString(detailFullText) {
				keyboard = "JIS"
			}
			if reUSKeyboard.MatchString(detailFullText) {
				keyboard = "US"
			}
		}

		products = append(products, &Product{
			ID:       productID,
			Name:     productName,
			Specs:    specs,
			Keyboard: keyboard,
			Price:    price,
			URL:      url,
		})
	})

	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonBytes)
}
