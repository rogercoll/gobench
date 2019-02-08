package gobench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

//To test them independentlly git test -run /POST
//Single benchmark go test -bench=BenchmarkLoad32 -run=^a

func TestHandler(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		valid := false
		req, err := http.NewRequest("GET", "/add", nil)
		if err != nil {
			t.Fatal(err)
			return
		}
		res := httptest.NewRecorder()
		NewServer().ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			valid = true
		} else {
			t.Fatalf("Bad status code, it cannot be a StatusOK code %d\n", res.Code)
		}
		assert.Equal(t, true, valid, "The two results should be the same.")

	})
	t.Run("POST", func(t *testing.T) {
		values := Values{
			A: 2,
			B: 3,
		}
		jsonData, _ := json.Marshal(values)
		req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		NewServer().ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("Bad status code %d\n", res.Code)
		}
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		bodyString := string(bodyBytes)
		goodResult := strconv.Itoa(values.A + values.B)
		assert.Equal(t, goodResult, bodyString, "The two results should be the same.")
	})
}

func request(router *mux.Router, client *http.Client, n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		req, err := http.NewRequest("GET", "/loadtest", nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		req.Header.Set("Connection", "close")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		io.Copy(ioutil.Discard, res.Body)
		if res.Code >= 200 && res.Code <= 299 {
			continue
		} else {
			log.Printf("HTTP Code = %d\n", res.Code)
		}
	}
	wg.Done()
}

//We make it private, avoids the testing driver trying to invoke it directly
func benchmarkLoad(c int, b *testing.B) {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = c
	defaultTransport.MaxIdleConnsPerHost = c
	client := &http.Client{
		Transport: &defaultTransport,
	}
	b.N = 10000
	var wg sync.WaitGroup
	router := NewServer()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go request(router, client, b.N, &wg)
	}
	wg.Wait()
}

func BenchmarkLoad1(b *testing.B)   { benchmarkLoad(1, b) }
func BenchmarkLoad8(b *testing.B)   { benchmarkLoad(8, b) }
func BenchmarkLoad16(b *testing.B)  { benchmarkLoad(16, b) }
func BenchmarkLoad32(b *testing.B)  { benchmarkLoad(32, b) }
func BenchmarkLoad64(b *testing.B)  { benchmarkLoad(64, b) }
func BenchmarkLoad128(b *testing.B) { benchmarkLoad(128, b) }
