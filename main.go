package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Specs []string `json:"specs"`
	Price string   `json:"price"`
}

func main() {
	doc, err := goquery.NewDocument("https://www.apple.com/jp/shop/browse/home/specialdeals/mac/macbook_pro/13")
	if err != nil {
		panic(err)
	}

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

		products = append(products, &Product{
			ID:    productID,
			Name:  productName,
			Specs: specs,
			Price: price,
		})
	})

	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonBytes)
}
