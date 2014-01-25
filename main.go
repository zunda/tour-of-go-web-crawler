package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Printer struct {
	ch chan bool
}

func (p *Printer) Print(s string) {
	if p.ch == nil {
		p.ch = make(chan bool)
		go func() { p.ch <- true }()
	}
	<-p.ch
	fmt.Print(s)
	go func() { p.ch <- true }()
}

type urlRecord struct {
	urls map[string]bool
	ch   chan bool
}

func (r *urlRecord) Init() {
	r.urls = map[string]bool{}
	r.ch = make(chan bool)
	go func() { r.ch <- true }()
}

func (r *urlRecord) Fetching(url string) bool {
	<-r.ch
	_, found := r.urls[url]
	if !found {
		r.urls[url] = true
	}
	go func() { r.ch <- true }()
	return !found
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, record *urlRecord, printer *Printer, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 {
		return
	}
	if record.Fetching(url) {
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			printer.Print(err.Error() + "\n")
			return
		}
		printer.Print(fmt.Sprintf("found: %s %q\n", url, body))
		for _, u := range urls {
			wg.Add(1)
			go Crawl(u, depth-1, fetcher, record, printer, wg)
		}
	} else {
		printer.Print(fmt.Sprintf("already crawled: %s\n", url))
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var printer Printer
	var urlrecord urlRecord
	urlrecord.Init()
	wg.Add(1)
	go Crawl("http://golang.org/", 4, fetcher, &urlrecord, &printer, &wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
