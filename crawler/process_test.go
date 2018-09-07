package crawler

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type openCloseBuffer struct {
	buffer *bytes.Buffer
}

func (b *openCloseBuffer) Read(p []byte) (n int, err error) {
	return b.buffer.Read(p)
}

func (b *openCloseBuffer) Close() error {
	return nil
}

func Test_ProcessPage_Simple(t *testing.T) {
	Convey("Given a simple html page", t, func() {
		newGetter := func(uri string) (*http.Response, error) {
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

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and create a basic page struct", func() {

			u, _ := url.Parse("http://local.link/zzzz")
			page, err := doProcessPage(u, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a title")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})
}

func Test_ProcessPage_RemoteLinks(t *testing.T) {
	Convey("Given an html page with remote links", t, func() {

		newGetter := func(uri string) (*http.Response, error) {
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

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and confirm the basic page struct has the remote pages", func() {
			d, _ := url.Parse("http://local.link")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a new article")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link")
			So(len(page.RemotePages), ShouldEqual, 2)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})
}

func Test_ProcessPage_LocalLinks(t *testing.T) {
	Convey("Given an html page with local links", t, func() {

		newGetter := func(uri string) (*http.Response, error) {
			var page string
			if uri == "http://local.link" {
				page = `
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
			} else if uri == "http://local.link/somewhere" {
				page = `
												<html>
																<head>
																				<title>This is a sub-article</title>
																</head>
																<body>
																				<p>This has no links</p>
																</body>
												</html>
								`
			} else if uri == "http://local.link/somewhere2" {
				page = `
												<html>
																<head>
																				<title>This is a sub-article #2</title>
																</head>
																<body>
																				<p>This has no links either</p>
																</body>
												</html>
								`

			} else {
				return nil, errors.New("Invalid url")
			}

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and confirm the basic page struct has the local pages", func() {
			d, _ := url.Parse("http://local.link")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a new new article")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 2)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})
}

func Test_ProcessPage_Images(t *testing.T) {
	Convey("Given an html page with images", t, func() {
		newGetter := func(uri string) (*http.Response, error) {
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

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and confirm the basic page struct has the images", func() {
			d, _ := url.Parse("http://local.link")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is an article with images")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 2)
			So(page.Assets[0].Type, ShouldEqual, AssetTypeIMG)
			So(page.Assets[0].URI, ShouldEqual, "image.jpg")
			So(page.Assets[1].Type, ShouldEqual, AssetTypeIMG)
			So(page.Assets[1].URI, ShouldEqual, "image2.jpg")
		})
	})
}

func Test_ProcessPage_CSS(t *testing.T) {
	Convey("Given an html page with css links", t, func() {
		newGetter := func(uri string) (*http.Response, error) {
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

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and confirm the basic page struct has the css", func() {
			d, _ := url.Parse("http://local.link")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is an article with links")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 2)
			So(page.Assets[0].Type, ShouldEqual, AssetTypeCSS)
			So(page.Assets[0].URI, ShouldEqual, "stylesheet1.css")
			So(page.Assets[1].Type, ShouldEqual, AssetTypeCSS)
			So(page.Assets[1].URI, ShouldEqual, "stylesheet2.css")
		})
	})
}

func Test_ProcessPage_JS(t *testing.T) {
	//TODO: This test is failing for some reason.
	//			If I comment out the first script it will find the second but not
	//			otherwise. I suspect that this might be a bug in the library but
	//			I will investigate later.
	Convey("Given an html page with external JS scripts", t, func() {

		newGetter := func(uri string) (*http.Response, error) {
			page := `
								<html>
												<head>
																<title>This is an article with JS</title>
												</head>
												<body>
																<script src="javascript.js" type="application/javascript"></script>
																<p>This is an article with JS</p>
																<script src="javascript2.js" type="text/javascript"></script>
												</body>
								</html>
				`

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and confirm the basic page struct has the javascript", func() {
			d, _ := url.Parse("http://local.link")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is an article with JS")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 2)
			So(page.Assets[0].Type, ShouldEqual, AssetTypeJS)
			So(page.Assets[0].URI, ShouldEqual, "javascript.js")
			So(page.Assets[1].Type, ShouldEqual, AssetTypeJS)
			So(page.Assets[1].URI, ShouldEqual, "javascript2.js")
		})
	})
}

func Test_ProcessPage_InfiniteLoopOfLocalPages_Single(t *testing.T) {
	Convey("Given an html page with one link to itself", t, func() {

		newGetter := func(uri string) (*http.Response, error) {
			page := `
								<html>
												<head>
																<title>This is a title</title>
												</head>
												<body>
																<h1>This is a title</h1>
																<a href="zzzz">Link to myself</a>
																<a href="zzzz#p1">Link to myself with fragment</a>
																<a href="javascript:doSomething();">Do some javascript</a>
												</body>
								</html>
				`

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and check that an infinite loop is not created", func() {
			d, _ := url.Parse("http://local.link/zzzz")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a title")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 0)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})
}

func Test_ProcessPage_InfiniteLoopOfLocalPages_Separation(t *testing.T) {
	Convey("Given an html page with a link to another page which then links to the first page", t, func() {
		newGetter := func(uri string) (*http.Response, error) {
			var page string
			if uri == "http://local.link/zzzz" {
				page = `
								<html>
												<head>
																<title>This is a title</title>
												</head>
												<body>
																<h1>This is a title</h1>
																<a href="yyyy">Link to another page</a>
												</body>
								</html>
				`
			} else if uri == "http://local.link/yyyy" {
				page = `
												<html>
																<head>
																				<title>This is a sub-article</title>
																</head>
																<body>
																				<a href="zzzz">Link to the first page</a>
																</body>
												</html>
								`
			} else {
				return nil, errors.New("Url invalid")
			}

			resp := new(http.Response)
			resp.StatusCode = 200
			resp.Body = &openCloseBuffer{bytes.NewBufferString(page)}

			return resp, nil
		}

		Convey("Process the page and check that an infinite loop is not created", func() {
			d, _ := url.Parse("http://local.link/zzzz")
			page, err := doProcessPage(d, newGetter, 100)
			So(err, ShouldBeNil)
			So(page, ShouldNotBeNil)
			So(page.Title, ShouldEqual, "This is a title")
			So(page.Type, ShouldEqual, AssetTypeHTML)
			So(page.URI.String(), ShouldEqual, "http://local.link/zzzz")
			So(len(page.RemotePages), ShouldEqual, 0)
			So(len(page.Pages), ShouldEqual, 1)
			So(len(page.Assets), ShouldEqual, 0)
		})
	})
}
