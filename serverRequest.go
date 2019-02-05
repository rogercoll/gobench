package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const n = 1000
const c = 6

var wg sync.WaitGroup

func request(url string) {
	st := time.Now()
	//http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	for i := 0; i < n; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("%s\n", err)
		}
		resp.Body.Close()
		code := resp.StatusCode
		if code >= 200 && code <= 299 {
			continue
		} else {
			log.Printf("HTTP Code = %d\n", code)
		}

	}
	f := time.Since(st).Seconds()
	fmt.Printf("Time for %d requests took up %.2fs\n", n, f)
	wg.Done()
}

func main() {
	start := time.Now()

	for i := 0; i < c; i++ {
		wg.Add(1)
		go request("http://localhost:8080/")
	}
	wg.Wait()

	finish := time.Since(start).Seconds()
	fmt.Printf("Total time for n = %d and c = %d took up %.2fs\n", n, c, finish)
}
