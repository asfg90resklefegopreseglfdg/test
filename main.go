package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const limit = 2

func main() {
	urls := []string{"https://golang.org", "https://gist.github.com/scripter-v/45bf032a7f357b15e9d0a8a9087ec53c", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org"}

	total := 0
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	quotaCh := make(chan struct{}, limit)
	for _, url := range urls {
		wg.Add(1)
		go func(url string, wg *sync.WaitGroup, total *int, mu *sync.Mutex, quotaCh chan struct{}) {
			defer wg.Done()
			quotaCh <- struct{}{}
			time.Sleep(time.Second * 2)
			client := &http.Client{}
			resp, err := client.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			count := strings.Count(string(body), "Go")
			fmt.Printf("Count for %s: %d\r\n", url, count)
			mu.Lock()
			*total += count
			mu.Unlock()
			<-quotaCh
		}(url, wg, &total, mu, quotaCh)
	}
	wg.Wait()
	fmt.Printf("Total: %d", total)
}
