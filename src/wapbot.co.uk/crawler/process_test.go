package crawler

import (
	"bytes"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func Test_ProcessPage(t *testing.T) {
	Convey("Given a simple html page", t, func() {
		page := `
								<html>
												<head>
																<title>This is a title</title>
												</head>
												<body>
																<h1>This is a title</h1>
												</body>
								</html>
				`

		Convey("Process the page and create a basic page struct", func() {

			u, _ := url.Parse("http://local.link/zzzz")
			page, err := ProcessPage(nil, u, bytes.NewBufferString(page), nil)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a title")
			So(page.Type, ShouldEqual, AssetType_HTML)
			So(page.URI, ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})

	Convey("Given an html page with remote links", t, func() {
		page := `
								<html>
												<head>
																<title>This is a new article</title>
												</head>
												<body>
																<a href="http://remotelink/somewhere">This is a remote link somewhere</a>
																<a href="http://remotelink/somewhere2">As is this</a>
												</body>
								</html>
					`

		Convey("Process the page and confirm the basic page struct has the remote pages", func() {
			d, _ := url.Parse("http://local.link")
			u, _ := url.Parse("http://local.link/zzzz")
			page, err := ProcessPage(d, u, bytes.NewBufferString(page), nil)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a new article")
			So(page.Type, ShouldEqual, AssetType_HTML)
			So(page.URI, ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 2)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})

	Convey("Given an html page with local links", t, func() {
		page := `
								<html>
												<head>
																<title>This is a new new article</title>
												</head>
												<body>
																<a href="http://local.link/somewhere">This is a remote link somewhere</a>
																<a href="somewhere2">As is this</a>
												</body>
								</html>
					`

		newGetter := func(uri string) (*httpResponse, error) {
			if uri == "http://local.link/somewhere" {
				page := `
												<html>
																<head>
																				<title>This is a sub-article</title>
																</head>
																<body>
																				<p>This has no links</p>
																</body>
												</html>
								`

				return &httpResponse{200, bytes.NewBufferString(page)}, nil

			} else if uri == "http://local.link/somewhere2" {
				page := `
												<html>
																<head>
																				<title>This is a sub-article #2</title>
																</head>
																<body>
																				<p>This has no links either</p>
																</body>
												</html>
								`

				return &httpResponse{200, bytes.NewBufferString(page)}, nil
			}

			return nil, errors.New("Invalid url")
		}

		Convey("Process the page and confirm the basic page struct has the local pages", func() {
			d, _ := url.Parse("http://local.link")
			u, _ := url.Parse("http://local.link/zzzz")
			page, err := ProcessPage(d, u, bytes.NewBufferString(page), newGetter)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a new new article")
			So(page.Type, ShouldEqual, AssetType_HTML)
			So(page.URI, ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 2)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})

	Convey("Given an html page with images", t, func() {
		page := `
								<html>
												<head>
																<title>This is an article with images</title>
												</head>
												<body>
																<img src="image.jpg"/>
																<img src="image2.jpg"/>
												</body>
								</html>
					`

		Convey("Process the page and confirm the basic page struct has the images", func() {
			d, _ := url.Parse("http://local.link")
			u, _ := url.Parse("http://local.link/zzzz")
			page, err := ProcessPage(d, u, bytes.NewBufferString(page), nil)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is an article with images")
			So(page.Type, ShouldEqual, AssetType_HTML)
			So(page.URI, ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 2)
			So(page.Assets[0].Type, ShouldEqual, AssetType_IMG)
			So(page.Assets[0].URI, ShouldEqual, "image.jpg")
			So(page.Assets[1].Type, ShouldEqual, AssetType_IMG)
			So(page.Assets[1].URI, ShouldEqual, "image2.jpg")
		})
	})

	Convey("Given an html page with css links", t, func() {
		page := `
								<html>
												<head>
																<title>This is an article with links</title>
																<link href="stylesheet1.css" rel="stylesheet" />
																<link href="stylesheet2.css" rel="stylesheet" />
												</head>
												<body>
																<p>This is an article with css</p>
												</body>
								</html>
					`

		Convey("Process the page and confirm the basic page struct has the css", func() {
			d, _ := url.Parse("http://local.link")
			u, _ := url.Parse("http://local.link/zzzz")
			page, err := ProcessPage(d, u, bytes.NewBufferString(page), nil)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is an article with links")
			So(page.Type, ShouldEqual, AssetType_HTML)
			So(page.URI, ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 2)
			So(page.Assets[0].Type, ShouldEqual, AssetType_CSS)
			So(page.Assets[0].URI, ShouldEqual, "stylesheet1.css")
			So(page.Assets[1].Type, ShouldEqual, AssetType_CSS)
			So(page.Assets[1].URI, ShouldEqual, "stylesheet2.css")
		})
	})

	//TODO: Javascript Assets
	//TODO: Local pages with loops
}
