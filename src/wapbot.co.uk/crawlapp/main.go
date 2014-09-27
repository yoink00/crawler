package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime/pprof"
	"wapbot.co.uk/crawler"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
