Crawler
=======

This is a simple one-site web-crawler written in Go

Install
-------
Assuming you've setup a Go development environment as per the instructions on the [Go Website].

Then installing is as simple as:

  - `git clone https://github.com/yoink00/crawler.git`
  - `cd crawler`
  - `source ENV; # This sets up your GOPATH and PATH environments. If you don't want that then skip`
  - `go get -d ./...`
  - `cd src/wapbot.co.uk/crawlapp`
  - `go install`
  - 
  
Running
-------
The above will have installed a binary called `crawlapp` in the `$GOPATH/bin` folder.

The binary accepts two flags:

  - `-cpuprofile=out_file` which outputs pprof compatible profiling information to `out_file`
  - `-site=site_to_search` which is the site that should be crawled.
  - 
  
Versions
--------
There are two tags on the master branch:

  - `initial_working_version` which contains the first working version.
  - `goroutine_working_version` which contains a version which spawns goroutines to fetch and process pages and is between 2x and 3x quicker that the first version.



License
----

MIT

[Go Website]:https://golang.org/doc/code.html

