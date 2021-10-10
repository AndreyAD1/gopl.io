package main

import (
	"fmt"
	"sort"
)

type Racer struct {
	Name string
	BirthYear int
	TitleNumber int
	FerrariRacer bool
}

type SortingTable struct {
	table []Racer
	lessFuncs []func (v1, v2 Racer) bool
}

func (t SortingTable) Len() int {
	return len(t.table)
}

func (t SortingTable) Swap(i, j int) {
	t.table[j], t.table[i] = t.table[i], t.table[j]
}

func (t SortingTable) Less(i, j int) bool {
	for _, lessFunc := range(t.lessFuncs) {
		switch {
		case lessFunc(t.table[i], t.table[j]):
			return true
		case lessFunc(t.table[j], t.table[i]):
			return false
		}
	}
	return false
}


func main() {
	racers := []Racer{
		{"Michael Schumacher", 1969, 7, true},
		{"Lewis Hamilton", 1985, 7, false},
		{"Ayrton Senna", 1960, 3, false},
	}
	nameIsLessFunc := func(r1, r2 Racer) bool {return r1.Name < r2.Name}
	titleNumIsLess := func(r1, r2 Racer) bool {return r1.TitleNumber < r2.TitleNumber}
	// birthYearIsLessFunc := func(r1, r2 Racer) bool {return r1.BirthYear < r2.BirthYear}
	funcs := []func (v1, v2 Racer) bool {titleNumIsLess, nameIsLessFunc}
	racerSortingTable := SortingTable{racers, funcs}
	sort.Sort(racerSortingTable)
	fmt.Println(racers)
}