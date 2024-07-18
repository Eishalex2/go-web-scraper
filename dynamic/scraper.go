package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Product struct {
	url, image, name, price string
}

func main() {
	var products []Product

	// initializing a chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// navigate to the target web page and select the HTML elements of interest
	var nodes []*cdp.Node
	chromedp.Run(ctx,
		chromedp.Navigate("https://www.scrapingcourse.com/ecommerce"),
		chromedp.Nodes(".product", &nodes, chromedp.ByQueryAll),
	)

	// scraping data from each node
	var url, image, name, price string
	for _, node := range nodes {
		chromedp.Run(ctx,
			chromedp.AttributeValue("a", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("img", "src", &image, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text("h2", &name, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".price", &price, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		product := Product{}

		product.url = url
		product.image = image
		product.name = name
		product.price = price

		products = append(products, product)
	}

	// export logic
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
