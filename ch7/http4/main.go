// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 195.

// Http4 is an e-commerce server that registers the /list and /price
// endpoint by calling http.HandleFunc.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var htmlTableTemplate = template.Must(template.New("htmlTable").Parse(`
<h1>Shop Items</h1>
<table>
<tr style='text-align: left'>
  <th>Item</th>
  <th>Price</th>
</tr>
{{range .}}
<tr>
  <td>{{.Item}}</a></td>
  <td>{{.Price}}</td>
</tr>
{{end}}
</table>
`))

type Errno uintptr

type ItemWithPrice struct {
	Item string
	Price dollars
}

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!-main

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	var items []ItemWithPrice
	for item, price := range db {
		items = append(items, ItemWithPrice{item, price})
		// fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	if err := htmlTableTemplate.Execute(w, items); err != nil {
		log.Fatal(err)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	price := req.URL.Query().Get("price")
	if price == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "'price' is required query argument")
		return
	}
	uintPrice, err := strconv.ParseUint(price, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "'price' should be a positive integer")
		return
	}
	itemName := req.URL.Query().Get("item")
	if itemName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "'itemName' is required query argument")
		return
	}
	db[itemName] = dollars(uintPrice)
	fmt.Fprintf(w, "%s: %s\n", itemName, db[itemName])
}
