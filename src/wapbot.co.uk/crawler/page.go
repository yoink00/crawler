package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
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
	sync.Mutex

	// The title of the page. Usually from <title> tags.
	Title string
	// The assets contained within the page
	Assets      []*Asset
	Pages       []*Page
	RemotePages []*Asset
}

func (p *Page) AddAsset(a *Asset) {
	p.Lock()
	defer p.Unlock()

	p.Assets = append(p.Assets, a)
}

func (p *Page) AddPage(np *Page) {
	p.Lock()
	defer p.Unlock()

	p.Pages = append(p.Pages, np)
}

func (p *Page) AddRemotePage(rp *Asset) {
	p.Lock()
	defer p.Unlock()

	p.RemotePages = append(p.RemotePages, rp)
}

/**
 * Dump data about this page and all pages it links to.
 */
func (p *Page) Dump() {
	var buf bytes.Buffer
	p.DumpToBuffer(&buf)
	fmt.Println(buf.String())
}

/**
 * Dump to buffer. This dumps the output to a buffer.
 */
func (p *Page) DumpToBuffer(buf *bytes.Buffer) {
	visited := make(map[string]bool)
	p.DumpPage_Indent(buf, 0, visited)
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
func (p *Page) DumpPage_Indent(buf *bytes.Buffer, level int, visited map[string]bool) {
	visited[p.URI] = true

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
			if !visited[np.URI] {
				np.DumpPage_Indent(buf, level+1, visited)
			} else {
				fmt.Fprintf(buf, "%sTitle: %s (previously visited)\n", indent(level+1), np.Title)
				fmt.Fprintf(buf, "%sURI:   %s\n", indent(level+1), np.URI)
			}
		}
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
