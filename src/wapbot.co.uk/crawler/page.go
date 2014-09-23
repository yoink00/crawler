package crawler

import (
	"bytes"
	"errors"
	"fmt"
)

type AssetType int

const (
	AssetType_JS AssetType = iota
	AssetType_HTML
	AssetType_CSS
	AssetType_IMG
)

func getTypeString(at AssetType) string {
	if at == AssetType_JS {
		return "JS"
	} else if at == AssetType_HTML {
		return "HTML"
	} else if at == AssetType_CSS {
		return "CSS"
	} else if at == AssetType_IMG {
		return "Image"
	} else {
		return "Unknown"
	}
}

/**
 * This struct describes a general asset.
 */
type Asset struct {
	URI  string
	Type AssetType
}

/**
 * This struct describes a web-page and the
 * pertinent information within.
 */
type Page struct {
	// Composed of asset
	Asset

	// The title of the page. Usually from <title> tags.
	Title string
	// The assets contained within the page
	Assets      []*Asset
	Pages       []*Page
	RemotePages []*Asset
}

func (p *Page) AddAsset(a *Asset) {
	p.Assets = append(p.Assets, a)
}

func (p *Page) AddPage(np *Page) {
	p.Pages = append(p.Pages, p)
}

func (p *Page) AddRemotePage(rp *Asset) {
	p.RemotePages = append(p.RemotePages, rp)
}

/**
 * Dump data about this page and all pages it links to.
 */
func (p *Page) Dump() {
	p.DumpPage_Indent(0)
}

/**
 * Indentation helper function
 */
func indent(level int) string {
	var indent bytes.Buffer
	for i := 0; i < level; i++ {
		indent.WriteString(" ")
	}

	return indent.String()
}

/**
 * Core dump page function with included indentation
 */
func (p *Page) DumpPage_Indent(level int) {
	fmt.Printf("%sTitle: %s\n", indent(level), p.Title)
	fmt.Printf("%sURI:   %s\n", indent(level), p.URI)
	fmt.Printf("%sAssets:\n", indent(level))

	for _, a := range p.Assets {
		fmt.Printf("%sURI: %s (%s)\n", indent(level+1), a.URI, getTypeString(a.Type))
	}

	fmt.Printf("%sRemote Pages:\n", indent(level))
	for _, rp := range p.RemotePages {
		fmt.Printf("%sURI: %s\n", indent(level+1), rp.URI)
	}

	fmt.Printf("%sPages:\n", indent(level))
	for _, np := range p.Pages {
		np.DumpPage_Indent(level + 1)
	}

	fmt.Println()
}

/**
 * Create a new asset struct and return the pointer
 */
func NewAsset(uri string, t AssetType) (*Asset, error) {

	if t < AssetType_JS || t > AssetType_IMG {
		return nil, errors.New("Invalid asset type")
	}

	asset := new(Asset)

	asset.URI = uri
	asset.Type = t

	return asset, nil
}

/**
 * Create a new page struct and return the pointer
 */
func NewPage(uri string, title string) *Page {
	page := new(Page)

	page.URI = uri
	page.Type = AssetType_HTML
	page.Title = title

	return page
}
