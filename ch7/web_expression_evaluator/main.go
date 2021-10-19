package main

import (
	"fmt"
	"log"
	"net/http"
)


func getMainPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive the request: %s %s", r.URL, r.Method)
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, MainPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrorPage)
	}
}

func getExpressionResult(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive the request: %s %s", r.URL, r.Method)
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, "WAT?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrorPage)
	}
}

func main() {
	http.HandleFunc("/", getMainPage)
	http.HandleFunc("/calculate", getExpressionResult)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}