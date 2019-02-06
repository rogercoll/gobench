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

//numb of conn per user/socket
const n = 10000

//numb of clients
const c = 2

var wg sync.WaitGroup

func request(url string, client *http.Client) {
	st := time.Now()
	for i := 0; i < n; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(ioutil.Discard, resp.Body)
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
	//Tenim 8 cores, a partir de c > 8 es reparteix les goroutines entre els 8, una goroutine a causa del Round Robin pot ser que passi de RUN a WAIT per això el temps no és igual a una quan c = 8
	start := time.Now()
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = c
	defaultTransport.MaxIdleConnsPerHost = c
	client := &http.Client{Transport: &defaultTransport}
	for i := 0; i < c; i++ {
		wg.Add(1)
		go request("http://localhost:8080/", client)
	}
	wg.Wait()

	finish := time.Since(start).Seconds()
	fmt.Printf("Total time for n = %d and c = %d took up %.2fs\n", n, c, finish)
}
