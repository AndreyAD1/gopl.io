// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package intset

import (
	"fmt"
	"testing"
)

func TestExampleOne(t *testing.T) {
	//!+main
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"
	//!-main

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
	if x.Len() != 4 {
		t.Errorf("Invalid length: expect: %v, got: %v", 4, x.Len())
	}
	x.Remove(9)
	if x.Has(9){
		t.Errorf("Got removed bit: %v", 9)
	}
	fmt.Println(x.String())
}

func TestExampleTwo(t *testing.T) {
	var x IntSet
	fmt.Println(x)
	x.Add(1)
	x.Add(144)
	x.Add(9)
	y := x.Copy()
	x.Add(42)

	//!+note
	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"
	//!-note

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
	x.Clear()
	fmt.Println(x)
	fmt.Println(y)
}
