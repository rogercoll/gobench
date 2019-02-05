package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	resp, err := http.Get("http://google.com")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("Correct fetch and response")
	} else {
		fmt.Println("Bad StatusCode")
	}
	fmt.Printf("%.2fs des del començament", time.Since(start).Seconds())
}
