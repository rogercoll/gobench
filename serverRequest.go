package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const n = 10000
const c = 9

var wg sync.WaitGroup

func request(url string) {
	st := time.Now()
	tr := &http.Transport{
		MaxIdleConns:       0,
		IdleConnTimeout:    0,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	for i := 0; i < n; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		//Setting the header to close is used to inform the server that the client wants to close the connection after the transaction is complete
		//Before doing that the connections were ESTABLISHED not in CLOSED state
		req.Header.Set("Connection", "close")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		} //if there isnt a continue then it prints the error and then dereferences a nil pointer(Body)
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
