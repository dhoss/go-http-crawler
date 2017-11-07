package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeMap struct {
	M map[string]bool
	L sync.Mutex
}

type Crawler struct {
	Listeners int
	Cond      *sync.Cond
}

var seen = &SafeMap{M: make(map[string]bool)}

var processed = make(chan string, 1)

func (m *SafeMap) Add(k string) {
	m.L.Lock()
	m.M[k] = true
	m.L.Unlock()
}

func (m *SafeMap) Get(k string) bool {
	m.L.Lock()
	defer m.L.Unlock()
	return m.M[k]
}

func Crawl(urls chan string) {
	for {
		for url := range urls {
			if seen.Get(url) {
				return
			}
			go fetch(url)
		}
	}

}

func fetch(url string) {
	seen.Add(url)
	processed <- url
}

func main() {
	var urls = make(chan string, 3)
	go Crawl(urls)
	urls <- "http://google.com"
	urls <- "http://bing.com"
	urls <- "http://facebook.com"
	urls <- "twitter.com"

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	for {

		select {
		case url := <-processed:
			// a read from ch has occurred
			fmt.Println("Processed ", url)
		case <-timeout:
			// the read from ch has timed out
			fmt.Println("Timeout")
			return
		}
	}

}
