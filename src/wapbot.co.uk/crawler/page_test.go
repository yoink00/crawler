package crawler

import (
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
