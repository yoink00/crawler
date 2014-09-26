package crawler

import (
	"fmt"
	"github.com/puerkitobio/goquery"
	"io"
	"net/url"
)

// A simple implementation of response to allow mocking in tests
type httpResponse struct {
	StatusCode int
	Body       io.Reader
}

type httpGetFunction func(string) (*httpResponse, error)

/*
 * Define a remote URL as one where:
 *	- The 'uri' is absolute
 *	- The 'domain' host and the 'uri' host do not match,
 *	- Or, if the hosts do match the paths do not match.
 */
func isRemoteLink(domain *url.URL, uri *url.URL) bool {
	if !uri.IsAbs() {
		fmt.Println(uri)
		uri = domain.ResolveReference(uri)
		fmt.Println(uri)
	}
	if uri.Host != domain.Host || (uri.Host != domain.Host && uri.Path != domain.Path) {
		return true
	}

	return false
}

func ProcessPage(domain *url.URL, uri *url.URL, buf io.Reader, getter httpGetFunction) (*Page, error) {
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return nil, err
	}

	title := doc.Find("title").Text()
	page := NewPage(uri.String(), title)

	doc.Find("a").EachWithBreak(func(_ int, sel *goquery.Selection) bool {
		href, exists := sel.Attr("href")
		if exists {
			newuri, err := url.Parse(href)
			if err != nil {
				return false
			}

			if isRemoteLink(domain, newuri) {
				rpage, err := NewAsset(href, AssetType_HTML)
				if err == nil {
					page.AddRemotePage(rpage)
				}
			} else {
				if !newuri.IsAbs() {
					newuri = domain.ResolveReference(newuri)
				}
				resp, err := getter(newuri.String())
				if err != nil {
					return false
				}
				newpage, err := ProcessPage(domain, newuri, resp.Body, getter)
				if err != nil {
					return false
				} else {
					page.AddPage(newpage)
				}
			}
		}

		return true
	})
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

	return page, nil
}
