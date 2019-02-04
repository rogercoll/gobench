package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	resp, err := http.Get("localhost:8000")
	if err != nil {
		fmt.Println("Error while requesting url, exit")
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error reading response")
		os.Exit(1)
	}
	fmt.Printf("%s", b)
	fmt.Printf("%.2fs des del començament", time.Since(start).Seconds())
}
