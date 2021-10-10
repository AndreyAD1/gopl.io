package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
)

type Racer struct {
	Name         string
	BirthYear    int
	TitleNumber  int
	FerrariRacer bool
}

type SortingTable struct {
	table     []Racer
	lessFuncs []func(v1, v2 Racer) bool
}

func (t SortingTable) Len() int {
	return len(t.table)
}

func (t SortingTable) Swap(i, j int) {
	t.table[j], t.table[i] = t.table[i], t.table[j]
}

func (t SortingTable) Less(i, j int) bool {
	for _, lessFunc := range t.lessFuncs {
		switch {
		case lessFunc(t.table[i], t.table[j]):
			return true
		case lessFunc(t.table[j], t.table[i]):
			return false
		}
	}
	return false
}

var htmlTableTemplate = template.Must(template.New("htmlTable").Parse(`
<h1>Racers</h1>
<table>
<tr style='text-align: left'>
  <th><a href="http://127.0.0.1:8000?sortby=name">Name</a></th>
  <th><a href="http://127.0.0.1:8000?sortby=birth_year">Birth Year</a></th>
  <th><a href="http://127.0.0.1:8000?sortby=title_number">Championships</a></th>
  <th><a href="http://127.0.0.1:8000?sortby=ferrari_racer">Ferrari Racer</a></th>
</tr>
{{range .}}
<tr>
  <td>{{.Name}}</a></td>
  <td>{{.BirthYear}}</td>
  <td>{{.TitleNumber}}</a></td>
  <td>{{.FerrariRacer}}</a></td>
</tr>
{{end}}
</table>
`))

var Racers = []Racer{
	{"Michael Schumacher", 1969, 7, true},
	{"Lewis Hamilton", 1985, 7, false},
	{"Ayrton Senna", 1960, 3, false},
	{"Juan Manuel Fangio", 1911, 5, true},
}

var funcsPerURLArgument = map[string]func(v1, v2 Racer) bool{
	"name":          func(r1, r2 Racer) bool { return r1.Name < r2.Name },
	"birth_year":    func(r1, r2 Racer) bool { return r1.BirthYear < r2.BirthYear },
	"title_number":  func(r1, r2 Racer) bool { return r1.TitleNumber < r2.TitleNumber },
	"ferrari_racer": func(r1, r2 Racer) bool { return r1.FerrariRacer && !r2.FerrariRacer },
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		getRacerTable(w, r)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func getRacerTable(out io.Writer, request *http.Request) {
	urlArgs := request.URL.Query()
	sortArg := urlArgs.Get("sortby")
	sortingFunc := funcsPerURLArgument[sortArg]

	if sortingFunc == nil {
		if err := htmlTableTemplate.Execute(out, Racers); err != nil {
			log.Fatal(err)
		}
		return
	}
	sortingFuncs := []func(v1, v2 Racer) bool{sortingFunc}
	racerSortingTable := SortingTable{Racers, sortingFuncs}
	sort.Sort(racerSortingTable)
	if err := htmlTableTemplate.Execute(out, Racers); err != nil {
		log.Fatal(err)
	}
}
