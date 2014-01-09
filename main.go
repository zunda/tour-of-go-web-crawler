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

func Print(ch chan string) {
	for s := range ch {
		fmt.Print(s)
	}
}

type urlRecord map[string]bool

func Fetching(url string, record chan urlRecord) (bool) {
	r := <- record
	_, found := r[url]
	if !found {
		r[url] = true
	}
	go func(){record <- r}()
	return !found
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, record chan urlRecord, printqueue chan string, wg *sync.WaitGroup) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	defer wg.Done()
	if depth <= 0 {
		return
	}
	if Fetching(url, record) {
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			printqueue <- err.Error() + "\n"
			return
		}
		printqueue <- fmt.Sprintf("found: %s %q\n", url, body)
		for _, u := range urls {
				wg.Add(1)
				go Crawl(u, depth-1, fetcher, record, printqueue, wg)
		}
	} else {
		printqueue <- fmt.Sprintf("already crawled: %s\n", url)
	}
	return
}

func main() {
	var wg sync.WaitGroup
	printqueue := make(chan string)
	record := make(chan urlRecord)
	go func(){record <- urlRecord{}}()
	go Print(printqueue)
	wg.Add(1)
	go Crawl("http://golang.org/", 4, fetcher, record, printqueue, &wg)
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
