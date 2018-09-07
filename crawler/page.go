package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
)

// AssetType describes the type of the asset on the page
type AssetType int

const (
	// AssetTypeJS is a JavaScript asset
	AssetTypeJS AssetType = iota
	// AssetTypeHTML is an HTML asset
	AssetTypeHTML
	// AssetTypeCSS is a style sheet asset
	AssetTypeCSS
	// AssetTypeIMG is an image asset
	AssetTypeIMG
)

func getTypeString(at AssetType) string {
	if at == AssetTypeJS {
		return "JS"
	} else if at == AssetTypeHTML {
		return "HTML"
	} else if at == AssetTypeCSS {
		return "CSS"
	} else if at == AssetTypeIMG {
		return "Image"
	} else {
		return "Unknown"
	}
}

// Asset describes a general asset.
type Asset struct {
	URI  string
	Type AssetType
}

// Page describes a web-page and the pertinent information within.
type Page struct {
	URI  *url.URL
	Type AssetType

	// The title of the page. Usually from <title> tags.
	Title string
	// The assets contained within the page
	Assets      []*Asset
	Pages       []*Page
	RemotePages []*Asset
}

// AddAsset adds the asset to the page
func (p *Page) AddAsset(a *Asset) {
	p.Assets = append(p.Assets, a)
}

// AddPage adds a local page (that will be later crawled) to the page
func (p *Page) AddPage(np *Page) {
	p.Pages = append(p.Pages, np)
}

// AddRemotePage adds a remote page (that will not be crawled) to the page
func (p *Page) AddRemotePage(rp *Asset) {
	p.RemotePages = append(p.RemotePages, rp)
}

// Dump will ouput data about this page and all pages it links to.
func (p *Page) Dump() {
	var buf bytes.Buffer
	p.DumpToBuffer(&buf)
	fmt.Println(buf.String())
}

// DumpToBuffer dumps the output to a buffer.
func (p *Page) DumpToBuffer(buf *bytes.Buffer) {
	visited := make(map[string]bool)
	p.DumpPageIndent(buf, 0, visited)
}

// Indentation helper function
func indent(level int) string {
	var indent bytes.Buffer
	for i := 0; i < level; i++ {
		indent.WriteString(" ")
	}

	return indent.String()
}

// DumpPageIndent dumps page function with included indentation
func (p *Page) DumpPageIndent(buf *bytes.Buffer, level int, visited map[string]bool) {
	visited[p.URI.String()] = true

	fmt.Fprintf(buf, "%sTitle: %s\n", indent(level), p.Title)
	fmt.Fprintf(buf, "%sURI:   %s\n", indent(level), p.URI)
	if len(p.Assets) > 0 {
		fmt.Fprintf(buf, "%sAssets:\n", indent(level))

		for _, a := range p.Assets {
			fmt.Fprintf(buf, "%sURI: %s (%s)\n", indent(level+1), a.URI, getTypeString(a.Type))
		}
	}

	if len(p.RemotePages) > 0 {
		fmt.Fprintf(buf, "%sRemote Pages:\n", indent(level))
		for _, rp := range p.RemotePages {
			fmt.Fprintf(buf, "%sURI: %s\n", indent(level+1), rp.URI)
		}
	}

	if len(p.Pages) > 0 {
		fmt.Fprintf(buf, "%sPages:\n", indent(level))
		for _, np := range p.Pages {
			if !visited[np.URI.String()] {
				np.DumpPageIndent(buf, level+1, visited)
			} else {
				fmt.Fprintf(buf, "%sTitle: %s (previously visited)\n", indent(level+1), np.Title)
				fmt.Fprintf(buf, "%sURI:   %s\n", indent(level+1), np.URI)
			}
		}
	}

	fmt.Println()
}

// NewAsset creates a new asset struct and returns the pointer
func NewAsset(uri string, t AssetType) (*Asset, error) {

	if t < AssetTypeJS || t > AssetTypeIMG {
		return nil, errors.New("Invalid asset type")
	}

	asset := new(Asset)

	asset.URI = uri
	asset.Type = t

	return asset, nil
}

// NewPage creates a new page struct and returns the pointer
func NewPage(uri *url.URL) *Page {
	page := new(Page)

	page.URI = uri
	page.Type = AssetTypeHTML

	return page
}
