package crawler

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_NewAsset(t *testing.T) {
	Convey("Create a JS asset", t, func() {
		asset, err := NewAsset("aaaa", AssetType_JS)
		So(err, ShouldBeNil)
		So(asset.URI, ShouldEqual, "aaaa")
		So(asset.Type, ShouldEqual, AssetType_JS)
	})
	Convey("Create a HTML asset", t, func() {
		asset, err := NewAsset("aaaa", AssetType_HTML)
		So(err, ShouldBeNil)
		So(asset.URI, ShouldEqual, "aaaa")
		So(asset.Type, ShouldEqual, AssetType_HTML)
	})
	Convey("Create a CSS asset", t, func() {
		asset, err := NewAsset("aaaa", AssetType_CSS)
		So(err, ShouldBeNil)
		So(asset.URI, ShouldEqual, "aaaa")
		So(asset.Type, ShouldEqual, AssetType_CSS)
	})
	Convey("Create a IMG asset", t, func() {
		asset, err := NewAsset("aaaa", AssetType_IMG)
		So(err, ShouldBeNil)
		So(asset.URI, ShouldEqual, "aaaa")
		So(asset.Type, ShouldEqual, AssetType_IMG)
	})

	Convey("Create an asset with an invalid type", t, func() {
		asset, err := NewAsset("aaaa", 9999)
		So(err, ShouldNotBeNil)
		So(asset, ShouldBeNil)
	})
}

func Test_NewPage(t *testing.T) {
	Convey("Create a new page", t, func() {
		page := NewPage("aaaa", "This is the title")
		So(page.Title, ShouldEqual, "This is the title")
		So(page.Type, ShouldEqual, AssetType_HTML)
		So(page.URI, ShouldEqual, "aaaa")
	})
}

func Test_DumpPage_Indent(t *testing.T) {
	Convey("Given a simple one page struct with values set", t, func() {
		page := NewPage("aaaa", "Title")

		Convey("Check that it is dumped correctly", func() {
			var buf bytes.Buffer
			page.DumpToBuffer(&buf)
			So(buf.String(), ShouldEqual, `Title: Title
URI:   aaaa
`)
		})

		Convey("Add some assets to the page", func() {
			asset1, err := NewAsset("bbbb.js", AssetType_JS)
			So(err, ShouldBeNil)
			asset2, err := NewAsset("cccc.js", AssetType_JS)
			So(err, ShouldBeNil)

			page.AddAsset(asset1)
			page.AddAsset(asset2)

			Convey("And check that the page with assets is dumped correctly", func() {
				var buf bytes.Buffer
				page.DumpToBuffer(&buf)
				So(buf.String(), ShouldEqual, `Title: Title
URI:   aaaa
Assets:
 URI: bbbb.js (JS)
 URI: cccc.js (JS)
`)
			})

			Convey("Add some remote pages", func() {
				rpage1, err := NewAsset("dddd", AssetType_HTML)
				So(err, ShouldBeNil)
				rpage2, err := NewAsset("eeee", AssetType_HTML)
				So(err, ShouldBeNil)

				page.AddRemotePage(rpage1)
				page.AddRemotePage(rpage2)

				Convey("And check that the page with remote pages is dumped correctly", func() {
					var buf bytes.Buffer
					page.DumpToBuffer(&buf)
					So(buf.String(), ShouldEqual, `Title: Title
URI:   aaaa
Assets:
 URI: bbbb.js (JS)
 URI: cccc.js (JS)
Remote Pages:
 URI: dddd
 URI: eeee
`)
				})

				Convey("Add a local page", func() {
					lpage := NewPage("ffff", "Title2")
					So(err, ShouldBeNil)

					page.AddPage(lpage)

					Convey("And check that the page with a local page is dumped correctly", func() {
						var buf bytes.Buffer
						page.DumpToBuffer(&buf)
						So(buf.String(), ShouldEqual, `Title: Title
URI:   aaaa
Assets:
 URI: bbbb.js (JS)
 URI: cccc.js (JS)
Remote Pages:
 URI: dddd
 URI: eeee
Pages:
 Title: Title2
 URI:   ffff
`)
					})

					//This tests that we do not end up in an infinite loop if we have already visited a node.
					//We will not print out the full details of a previously visited node but just a summary.
					//TODO: Make this detect if we are in a loop (by for example breaking out after a few seconds)
					//      At the moment the only way to know the test has failed is by the process eating up all
					//      our memory and taking a long time.
					Convey("Add a reference to the new page to the original page to create a loop", func() {
						lpage.AddPage(page)

						Convey("And test that we do not end up in an infinite loop", func() {

							var buf bytes.Buffer
							page.DumpToBuffer(&buf)
							So(buf.String(), ShouldEqual, `Title: Title
URI:   aaaa
Assets:
 URI: bbbb.js (JS)
 URI: cccc.js (JS)
Remote Pages:
 URI: dddd
 URI: eeee
Pages:
 Title: Title2
 URI:   ffff
 Pages:
  Title: Title (previously visited)
  URI:   aaaa
`)
						})
					})
				})
			})
		})
	})
}
