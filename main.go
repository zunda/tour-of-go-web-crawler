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
	m sync.Mutex
}

func (p *Printer) Print(s string) {
	p.m.Lock()
	fmt.Print(s)
	p.m.Unlock()
}

// urlRecord records URLs that have been fetched
type urlRecord struct {
	urls map[string]bool
	m    sync.Mutex
}

// returns true when the url is not yet fetched and records it as fetched
func (r *urlRecord) Fetching(url string) bool {
	if r.urls == nil {
		r.urls = map[string]bool{}
	}
	r.m.Lock()
	_, found := r.urls[url]
	if !found {
		r.urls[url] = true
	}
	r.m.Unlock()
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
