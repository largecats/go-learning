package main

import (
	"fmt"
	"sync"
)

// SafeMap is safe to use concurrently.
type SafeMap struct {
	mu sync.Mutex
	v  map[string]string
}

// Set sets value for given key
func (c *SafeMap) Set(key string, value string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = value
	c.mu.Unlock()
}

// Value returns mapped value for given key
func (c *SafeMap) Value(key string) (string, bool) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	val, ok := c.v[key]
	return val, ok
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, urlChannel chan string, urlMap SafeMap) {
	defer close(urlChannel) // close channel right before Crawl returns
	if depth <= 0 {
		return
	}

	_, ok := urlMap.Value(url)
	if !ok { // don't visit the same url twice
		body, urls, err := fetcher.Fetch(url)
		urlMap.Set(url, body)

		if err != nil {
			urlChannel <- err.Error()
			return
		}

		urlChannel <- fmt.Sprintf("found: %s %q", url, body)

		subPagesChannel := make([]chan string, len(urls)) // store results for sub-pages on given url
		for i, u := range urls {
			subPagesChannel[i] = make(chan string)
			go Crawl(u, depth-1, fetcher, subPagesChannel[i], urlMap)
		}

		for i := range subPagesChannel {
			for s := range subPagesChannel[i] {
				urlChannel <- s
			}
		}
	}

	return
}

func main() {
	urlChannel := make(chan string)
	urlMap := SafeMap{v: make(map[string]string)}
	go Crawl("https://golang.org/", 4, fetcher, urlChannel, urlMap)

	for s := range urlChannel {
		fmt.Println(s)
	}
}

// fakeFetcher is Fetcher that returns canned urlChannels.
type fakeFetcher map[string]*fakeurlChannel

type fakeurlChannel struct {
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
	"https://golang.org/": &fakeurlChannel{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeurlChannel{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeurlChannel{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeurlChannel{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

/*
found: https://golang.org/ "The Go Programming Language"
found: https://golang.org/pkg/ "Packages"
found: https://golang.org/pkg/fmt/ "Package fmt"
found: https://golang.org/pkg/os/ "Package os"
not found: https://golang.org/cmd/
*/
