package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"

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

func fetchProducts() ([]*Product, error) {
	doc, err := goquery.NewDocument("https://www.apple.com/jp/shop/browse/home/specialdeals/mac/macbook_pro/13")
	if err != nil {
		return nil, err
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

		url := specsSelection.Find("h3 > a").AttrOr("href", "#")
		if url != "#" {
			url = "https://www.apple.com" + url
		}

		products = append(products, &Product{
			ID:       productID,
			Name:     productName,
			Specs:    specs,
			Keyboard: "UNKNOWN",
			Price:    price,
			URL:      url,
		})
	})

	return products, nil
}

func (p *Product) fetchDetails() {
	if p.URL != "#" {
		detailDoc, err := goquery.NewDocument(p.URL)
		if err != nil {
			p.Keyboard = "ERROR"
		}

		detailFullText := detailDoc.Text()
		if regexp.MustCompile(`JIS.*キーボード`).MatchString(detailFullText) {
			p.Keyboard = "JIS"
		}
		if regexp.MustCompile(`US.*キーボード`).MatchString(detailFullText) {
			p.Keyboard = "US"
		}
	}
}

func main() {
	threads := flag.Int("threads", 5, "")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	products, err := fetchProducts()
	if err != nil {
		panic(err)
	}

	// make worker threads
	queue := make(chan *Product, len(products))
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				product, ok := <-queue
				if !ok {
					return
				}
				product.fetchDetails()
			}
		}()
	}

	for _, product := range products {
		queue <- product
	}
	close(queue)
	wg.Wait()

	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonBytes)
}
