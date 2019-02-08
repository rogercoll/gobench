package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/add", nil)
		if err != nil {
			fmt.Println("Aqui arriba")
			t.Fatal(err)
			return
		}
		res := httptest.NewRecorder()
		NewServer().ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("Bad status code %d\n", res.Code)
		}
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
