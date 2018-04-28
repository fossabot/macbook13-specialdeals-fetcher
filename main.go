package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Release   string `json:"release"`
	Processor string `json:"processor"`
	Memory    string `json:"memory"`
	Storage   string `json:"storage"`
	Keyboard  string `json:"keyboard"`
	Price     string `json:"price"`
	URL       string `json:"url"`
}

type ProductParser struct {
	Document *goquery.Document
}

func (p *Product) LoadFromURL(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	parser := NewProductParser(doc)

	p.ID = parser.GetProductID()
	p.Name = parser.GetProductName()
	p.Release = parser.GetRelease()
	p.Processor = parser.GetProcessor()
	p.Memory = parser.GetMemory()
	p.Storage = parser.GetStorage()
	p.Keyboard = parser.GetKeyboard()
	p.Price = parser.GetPrice()

	p.URL = url

	return nil
}

func (pp *ProductParser) GetProductID() string {
	id, ok := pp.Document.Find("input[name=product]").Attr("value")
	if !ok {
		return "UNKNOWN"
	}

	return id
}

func (pp *ProductParser) GetProductName() string {
	return strings.TrimSpace(
		pp.Document.Find("h1[data-autom=productTitle]").Text(),
	)
}

func (pp *ProductParser) GetRelease() string {
	return strings.TrimSpace(
		pp.Document.Find(".Overview-panel .para-list").Eq(0).Text(),
	)
}

func (pp *ProductParser) GetProcessor() string {
	return strings.TrimSpace(pp.Document.Find("#techSpecsSection .para-list").Eq(0).Text())
}

func (pp *ProductParser) GetMemory() string {
	return strings.TrimSpace(pp.Document.Find("#techSpecsSection .para-list").Eq(1).Text())
}

func (pp *ProductParser) GetStorage() string {
	s := strings.TrimSpace(pp.Document.Find("#techSpecsSection .para-list").Eq(2).Text())
	return s[:len(s)-1]
}

func (pp *ProductParser) GetKeyboard() string {
	text := pp.Document.Text()

	if regexp.MustCompile(`(?i)JIS.*(キーボード|Key)`).MatchString(text) {
		return "JIS"
	}

	if regexp.MustCompile(`(?i)(US|U\.S\.).*(キーボード|Key)`).MatchString(text) {
		return "US"
	}

	return "UNKNOWN"
}

func (pp *ProductParser) GetPrice() string {
	return strings.TrimSpace(pp.Document.Find(".current_price").Text())
}

func NewProduct() *Product {
	return &Product{}
}

func NewProductParser(doc *goquery.Document) *ProductParser {
	return &ProductParser{
		Document: doc,
	}
}

func fetchProductURLs(locale string) ([]string, error) {
	doc, err := goquery.NewDocument("https://www.apple.com/" + locale + "/shop/browse/home/specialdeals/mac/macbook_pro/13")
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("tr.product").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Find("a.button").Attr("href")
		if !ok {
			return
		}
		url = "https://www.apple.com" + url

		urls = append(urls, url)
	})

	return urls, nil
}

func main() {
	threads := flag.Int("threads", 5, "")
	locale := flag.String("locale", "jp", "")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	urls, err := fetchProductURLs(*locale)
	if err != nil {
		panic(err)
	}

	// make worker threads
	var products []*Product
	queue := make(chan string, len(urls))
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				url, ok := <-queue
				if !ok {
					return
				}

				product := NewProduct()
				err := product.LoadFromURL(url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					continue
				}
				products = append(products, product)
			}
		}()
	}

	for _, url := range urls {
		queue <- url
	}
	close(queue)
	wg.Wait()

	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonBytes)
}
