package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

var c = colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
)

func main() {
	

	var visitedUrls sync.Map

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	
	err := c.SetProxy("http://35.185.196.38:3128")
	
	if err != nil {
		log.Fatal(err)
	}
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// });
	
	// c.OnError(func(_ *colly.Response, err error){
	// 	fmt.Println("Something went wrong",err)
	// });
	
	// c.OnResponse(func(r *colly.Response){
	// 	fmt.Println("Page Visited",r.Request.URL)
	// })
	
	// c.OnHTML("a", func(e *colly.HTMLElement){
	// 	fmt.Println("Link Found",e.Text,e.Attr("href"))
	// });
	
	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println(r.Request.URL,"scraped")
	// });
	
	// 
    var products []Product
    
   	c.OnHTML("li.product", func(e *colly.HTMLElement){
    
        product := Product{}
        product.Url = e.ChildAttr("a","href")
        product.Image = e.ChildAttr("img","src")
        product.Name = e.ChildText(".product-name")
        product.Price = e.ChildText(".price")
        products = append(products,product)
    
		fmt.Println("Link Found",e.Text,e.Attr("href"))
	});
    
    c.OnHTML("a.next", func(e *colly.HTMLElement){
		nextPage := e.Attr("href")
		if _,found := visitedUrls.Load(nextPage); !found {
			fmt.Println("scraping",nextPage)
			visitedUrls.Store(nextPage,true)
			e.Request.Visit(nextPage)
		}
	})
    
    c.OnScraped(func(r *colly.Response){
        file,err := os.Create("products.csv")
        if err != nil {
            log.Fatalln("failed to create csv")
        }
        defer file.Close()
        writer := csv.NewWriter(file)
        headers := []string { "Url", "Image", "Name", "Price" }
        writer.Write(headers)
        for _, product := range products {
            writer.Write([]string{product.Url, product.Image, product.Name, product.Price})
        }
        writer.Flush()
    })
    
   	 c.Visit("https://www.scrapingcourse.com/ecommerce")


	
}

