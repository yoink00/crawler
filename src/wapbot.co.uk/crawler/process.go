package crawler

import (
	"github.com/puerkitobio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type httpGetFunction func(string) (*http.Response, error)

/*
 * Define a remote URL as one where:
 *	- The 'uri' is absolute
 *	- The 'domain' host and the 'uri' host do not match,
 *	- Or, if the hosts do match the paths do not match.
 */
func isRemoteLink(domain *url.URL, uri *url.URL) bool {
	if !uri.IsAbs() {
		uri = domain.ResolveReference(uri)
	}
	if uri.Host != domain.Host || (uri.Host == domain.Host && !strings.HasPrefix(uri.Path, domain.Path)) {
		return true
	}

	return false
}

func isSameUri(domain *url.URL, uri *url.URL, newuri *url.URL) bool {
	if strings.HasPrefix(newuri.String(), "javascript") {
		return true
	}

	if !uri.IsAbs() {
		uri = domain.ResolveReference(uri)
		uri.Fragment = ""
	}

	if !newuri.IsAbs() {
		newuri = domain.ResolveReference(newuri)
		newuri.Fragment = ""
	}

	if uri.String() == newuri.String() {
		return true
	}

	return false
}

func doProcessPage(domain *url.URL, uri *url.URL, buf io.ReadCloser, getter httpGetFunction, visited map[string]*Page) (*Page, error) {
	defer buf.Close()

	// Check to see if we've visited this page.
	uri.Fragment = ""
	if visited != nil {
		if page, exists := visited[uri.String()]; exists {
			// If we have then return it.
			return page, nil
		} else if page, exists := visited[uri.String()+"/"]; exists {
			return page, nil
		} else if page, exists := visited[strings.TrimRight(uri.String(), "/")]; exists {
			return page, nil
		}
	}

	// Process the new document
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").Text()
	page := NewPage(uri.String(), title)

	if visited != nil {
		visited[uri.String()] = page
	}

	var wg sync.WaitGroup
	spawned := 0
	doc.Find("a").EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		href, exists := sel.Attr("href")
		if exists {
			newuri, err := url.Parse(href)
			if err != nil {
				return false
			}

			//If this is a link back to the same page then ignore it.
			if !isSameUri(domain, uri, newuri) {

				if isRemoteLink(domain, newuri) {
					rpage, err := NewAsset(href, AssetType_HTML)
					if err == nil {
						page.AddRemotePage(rpage)
					}
				} else {
					if !newuri.IsAbs() {
						newuri = domain.ResolveReference(newuri)
						newuri.Fragment = ""
					}
					if newpage, exists := visited[newuri.String()]; exists {
						page.AddPage(newpage)
					} else {
						wg.Add(1)
						spawned++
						go func() {
							defer wg.Done()
							resp, err := getter(newuri.String())
							if err != nil {
								return
							}
							if contentType, exists := resp.Header["Content-Type"]; exists {
								ok := false
								for _, s := range contentType {
									if strings.Contains(s, "text/html") {
										ok = true
									}
								}
								if !ok {
									return
								}
							}
							newpage, err := doProcessPage(domain, newuri, resp.Body, getter, visited)
							if err != nil {
								return
							} else {
								page.AddPage(newpage)
							}
						}()
					}
				}
			}
		}

		return true
	})
	wg.Wait()

	if err != nil {
		return nil, err
	}

	doc.Find("img").Each(func(_ int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if exists {
			asset, _ := NewAsset(src, AssetType_IMG)
			page.AddAsset(asset)
		}
	})

	doc.Find("link").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			rel, exists := sel.Attr("rel")
			if exists && rel == "stylesheet" {
				asset, _ := NewAsset(href, AssetType_CSS)
				page.AddAsset(asset)
			}
		}
	})

	doc.Find("script").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("src")
		if exists {
			typ, exists := sel.Attr("type")
			if exists &&
				(strings.HasSuffix(typ, "javascript") ||
					strings.HasSuffix(typ, "ecmascript")) {
				asset, _ := NewAsset(href, AssetType_JS)
				page.AddAsset(asset)
			}
		}
	})

	return page, nil
}

func ProcessPage(uri *url.URL) (*Page, error) {
	visited := make(map[string]*Page)

	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	return doProcessPage(uri, uri, resp.Body, http.Get, visited)
}
