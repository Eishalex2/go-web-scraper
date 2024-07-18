package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	url   string
	image string
	name  string
	price string
}

func main() {
	// initializing the slice of structs to store the data to scrape
	var products []Product
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := "https://www.scrapingcourse.com/ecommerce/page/1/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// iteration step
	i := 1
	limit := 5

	// creating a new Colly instance
	c := colly.NewCollector()

	// setting a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	// crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		newPaginationLink := e.Attr("href")

		// if the page is new
		if !slices.Contains(pagesToScrape, newPaginationLink) {
			// should page be scraped
			if !slices.Contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// scraping logic
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		product := Product{}

		product.url = e.ChildAttr("a", "href")
		product.image = e.ChildAttr("img", "src")
		product.name = e.ChildText("h2")
		product.price = e.ChildText(".price")

		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {
		// Printing which page is being scraped
		fmt.Println("Scraping page", i)

		// go while there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]
			i++

			// visit new page
			c.Visit(pageToScrape)
		} else {
			// opening the CSV file
			file, err := os.Create("products.csv")
			if err != nil {
				log.Fatalln("Failed to create output CSV file", err)
			}
			defer file.Close()

			// initializing a file writer
			writer := csv.NewWriter(file)

			// writing the CSV headers
			headers := []string{
				"url",
				"image",
				"name",
				"price",
			}
			writer.Write(headers)

			// writing each product as a CSV row
			for _, product := range products {
				// converting a Product to an array of strings
				record := []string{
					product.url,
					product.image,
					product.name,
					product.price,
				}

				// adding a CSV record to the output file
				writer.Write(record)
			}
			defer writer.Flush()
		}
	})

	// downloading the target HTML page
	c.Visit(pageToScrape)
}
