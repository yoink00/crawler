package main

import (
	"fmt"
	"net/url"
	"wapbot.co.uk/crawler"
)

func main() {

	uri, err := url.Parse("http://tomblomfield.com")
	if err != nil {
		fmt.Printf("Invalid url: %s\n", err.Error())
		return
	}

	page, err := crawler.ProcessPage(uri)
	if err != nil {
		fmt.Printf("Unable to crawl page: %s\n", err.Error())
		return
	}

	page.Dump()
}
