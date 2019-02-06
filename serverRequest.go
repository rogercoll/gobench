package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const n = 10000
const c = 5

var wg sync.WaitGroup

func request(url string) {
	st := time.Now()
	//http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000
	for i := 0; i < n; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		} //if there isnt a continue then it prints the error and then dereferences a nil pointer(Body)
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			continue
		} else {
			log.Printf("HTTP Code = %d\n", resp.StatusCode)
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
