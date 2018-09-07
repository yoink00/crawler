package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/yoink00/crawler/crawler"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var site = flag.String("site", "", "site to process")
var numPages = flag.Int("pages", 100, "number of pages to crawl")

func main() {
	flag.Parse()

	if *site == "" {
		fmt.Println("-site flag is mandatory")
		return
	}

	if !strings.HasPrefix(*site, "http://") &&
		!strings.HasPrefix(*site, "https://") {
		fmt.Println("-site must be a fully formed URL")
		return
	}

	fmt.Printf("GOMAXPROCS is set to: %d\n", runtime.GOMAXPROCS(-1))

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	uri, err := url.Parse(*site)
	if err != nil {
		fmt.Printf("Invalid url: %s\n", err.Error())
		return
	}

	page, err := crawler.ProcessPage(uri, numPages)
	if err != nil {
		fmt.Printf("Unable to crawl page: %s\n", err.Error())
		return
	}

	page.Dump()
}
