package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Values struct {
	A int `json:"a"`
	B int `json:"b"`
}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var m Values
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	sum := m.A + m.B
	fmt.Fprintf(w, "%d", sum)
}

func NewServer() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/add", handlerAdd)
	return router
}
