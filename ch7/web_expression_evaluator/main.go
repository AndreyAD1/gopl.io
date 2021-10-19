package main

import (
	"fmt"
	"log"
	"net/http"
)

type resultStruct struct {
	Expression string
	Result float64
}


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
	expression := r.URL.Query().Get("expression")
	w.Header().Set("Content-Type", "text/html")
	if expression == "" {
		result := resultStruct{"No expression", 0}
		w.WriteHeader(http.StatusBadRequest)
		err := ResultPageTemplate.Execute(w, result)
		if err != nil {
			log.Printf("Can not generate no expresion response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, ErrorPage)
		}
		return
	}
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