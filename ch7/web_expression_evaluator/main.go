package main

import (
	"log"
	"net/http"
)

func getMainPage(w http.ResponseWriter, r *http.Request) {}

func getExpressionResult(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/", getMainPage)
	http.HandleFunc("/calculate", getExpressionResult)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}