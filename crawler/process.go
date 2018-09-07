package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/puerkitobio/goquery"
)

type httpGetFunction func(string) (*http.Response, error)

/*
 * Define a remote URL as one where:
 *	- The 'domain' host and the 'uri' host do not match,
 */
func isRemoteLink(domain *url.URL, uri *url.URL) bool {
	if !uri.IsAbs() || !domain.IsAbs() {
		return true
	}

	if uri.Host != domain.Host {
		return true
	}

	return false
}

func isSameURI(uri *url.URL, newuri *url.URL) bool {
	if strings.HasPrefix(newuri.String(), "javascript") {
		log.Print("This is a JavaScript file: ", newuri.String())
		return true
	}

	if !uri.IsAbs() || !newuri.IsAbs() {
		log.Printf("%s or %s is not absolute", uri.String(), newuri.String())
		return true
	}

	if uri.String()[4:] == newuri.String()[4:] || uri.String()[5:] == newuri.String()[5:] {
		log.Println("The URIs match: ", uri, newuri)
		return true
	}

	return false
}

// processPage
// 1. Pull page object from process page channel
// 2. Get page
// 3. Extract assets and add to page
// 4. Extract new pages
// 5. Filter new pages into remote and local and add to page
// 6. Put parent page into post-process channel
// 7. Go to 1
func processPageFromChannel(pagesToVisit <-chan *Page, visitedPages chan<- *Page, getter httpGetFunction) {
	//defer close(visitedPages)

	for page := range pagesToVisit {
		// Process the new document
		log.Print("Visiting: ", page.URI.String())

		resp, err := getter(page.URI.String())
		if err != nil {
			log.Print("Error fetching page: ", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Print("Error processing page: ", err)
			continue
		}

		title := doc.Find("title").Text()
		page.Title = title

		doc.Find("a").Each(func(_ int, sel *goquery.Selection) {
			href, exists := sel.Attr("href")
			if exists {

				if strings.HasPrefix(href, "//") {
					// This is a absolute link but with no transport. We'll add https
					href = "https:" + href
				}
				newuri, err := page.URI.Parse(href)
				if err != nil {
					log.Print("Unable to parse URL: ", err)
				}
				newuri.Fragment = ""

				log.Print("Got URL from page: ", newuri.String())

				//If this is a link back to the same page then ignore it.
				if !isSameURI(page.URI, newuri) {
					log.Print("This is not a link to myself: ", newuri.String())

					if isRemoteLink(page.URI, newuri) {
						rpage, err := NewAsset(href, AssetTypeHTML)
						if err == nil {
							log.Print("Adding remote page: ", href)
							page.AddRemotePage(rpage)
						}
					} else {
						newpage := NewPage(newuri)
						page.AddPage(newpage)
					}
				} else {
					log.Printf("%s is a link back to %s", newuri, page.URI.String())
				}
			}
		})

		doc.Find("img").Each(func(_ int, sel *goquery.Selection) {
			src, exists := sel.Attr("src")
			if exists {
				asset, _ := NewAsset(src, AssetTypeIMG)
				page.AddAsset(asset)
			}
		})

		doc.Find("link").Each(func(_ int, sel *goquery.Selection) {
			href, exists := sel.Attr("href")
			if exists {
				rel, exists := sel.Attr("rel")
				if exists && rel == "stylesheet" {
					asset, _ := NewAsset(href, AssetTypeCSS)
					page.AddAsset(asset)
				}
			}
		})

		doc.Find("script").Each(func(_ int, sel *goquery.Selection) {
			href, exists := sel.Attr("src")
			if exists {
				log.Print("Got script: ", href)
				typ, exists := sel.Attr("type")
				if exists &&
					(strings.HasSuffix(typ, "javascript") ||
						strings.HasSuffix(typ, "ecmascript")) {
					asset, _ := NewAsset(href, AssetTypeJS)
					page.AddAsset(asset)
				}
			}
		})

		log.Printf("Adding page to %s visitedPages channel", page.URI.String())
		visitedPages <- page
	}
}

// postProcessPage
// 1. Pull page object from channel
// 2. Check if page is in visited map; if not add it
// 3. For each child local page:
//   a. Check if page is in visited map; if so swap with that; if not put into process page channel unless visited map >= 100 then close the channel
//
func postProcessPage(visitedPages <-chan *Page, pagesToVisit chan<- *Page, numPages int, wg *sync.WaitGroup) {
	defer func() {
		log.Print("Closing pagesToVisit channel")
		close(pagesToVisit)
		wg.Done()
	}()

	visited := make(map[string]*Page)

	msgReceived := false

loop:
	for {
		select {
		case page := <-visitedPages:
			msgReceived = true

			log.Print("Received new visited page: ", page.URI.String())

			if _, exists := visited[page.URI.String()]; !exists {
				log.Printf("%s has not already been visited", page.URI.String())
				visited[page.URI.String()] = page
			}

			if len(visited) > numPages {
				log.Printf("Visited %d pages", len(visited))
				break loop
			}

			for i := range page.Pages {
				if p, exists := visited[page.Pages[i].URI.String()]; exists {
					log.Printf("%s has already been visited", page.Pages[i].URI.String())
					page.Pages[i] = p
				} else {
					visited[page.Pages[i].URI.String()] = page.Pages[i]
					log.Printf("Visited %d pages", len(visited))
					pagesToVisit <- page.Pages[i]
				}
			}
		case <-time.Tick(1 * time.Second):
			log.Print("Timer fired")
			if msgReceived {
				log.Print("Resetting timer")
				msgReceived = false
			} else {
				log.Print("Timed out")
				break loop
			}
		}
	}
}

func doProcessPage(domain *url.URL, getter httpGetFunction, numPages int) (*Page, error) {
	page := NewPage(domain)

	pagesToVisit := make(chan *Page, 10000)
	visitedPages := make(chan *Page, 10000)

	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)
	go processPageFromChannel(pagesToVisit, visitedPages, getter)

	var wg sync.WaitGroup
	wg.Add(1)
	go postProcessPage(visitedPages, pagesToVisit, numPages, &wg)
	pagesToVisit <- page
	wg.Wait()

	return page, nil
}

// ProcessPage crawls a website starting at the specified URL. It returns the root page.
func ProcessPage(uri *url.URL, numPages int) (*Page, error) {
	return doProcessPage(uri, http.Get, numPages)
}
