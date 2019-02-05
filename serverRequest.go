package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const n = 100
const c = 15

var wg sync.WaitGroup

func request() {
	for i := 0; i < n; i++ {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			continue
		} else {
			fmt.Println("Bad StatusCode")
			os.Exit(1)
		}
	}
	defer wg.Done()
}

func main() {
	start := time.Now()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go request()
	}
	wg.Wait()
	finish := time.Since(start).Seconds()
	fmt.Printf("%.2fs des del començament\n", finish)
	fmt.Println("Everything OK!")
	s := fmt.Sprintf("For n = %d and c = %d took up %.2fs\n", n, c, finish)
	f, err := os.OpenFile("request_go_server.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(s)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
