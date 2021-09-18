// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"fmt"
)

const UintSize = 32 << (^uint(0) >> 63)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/UintSize, uint(x%UintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Elems() []int {
	var elements []int
	for wordIndex, word := range s.words {
		if word != 0 {
			for bitIndex := 0; bitIndex < UintSize; bitIndex++ {
				if word&(1<<bitIndex) != 0 {
					element := wordIndex*UintSize + bitIndex
					elements = append(elements, element)
				}
			}
		}
	}
	return elements
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/UintSize, uint(x%UintSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(otherSet *IntSet) {
	if len(s.words) > len(otherSet.words) {
		s.words = s.words[:len(otherSet.words)-1]
	}
	for i, otherSetWord := range otherSet.words {
		if i < len(s.words) {
			s.words[i] &= otherSetWord
		}
	}
}

func (s *IntSet) SymmetricDifference(otherSet *IntSet) {
	for i, otherSetWord := range otherSet.words {
		if i < len(s.words) {
			s.words[i] ^= otherSetWord
		} else {
			s.words = append(s.words, otherSetWord)
		}
	}
}

func (s *IntSet) DifferenceWith(otherSet *IntSet) {
	initialSet := s.Copy()
	s.SymmetricDifference(otherSet)
	s.IntersectWith(initialSet)
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < UintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", UintSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

func (s *IntSet) Len() int {
	var setBitsNumber int
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for bitIndex := 0; bitIndex < UintSize; bitIndex++ {
			if word&(1<<uint8(bitIndex)) != 0 {
				setBitsNumber++
			}
		}
	}
	return setBitsNumber
}

func (s *IntSet) Remove(x int) {
	word, bit := x/UintSize, uint(x%UintSize)
	s.words[word] &= ^(1 << bit)
}

func (s *IntSet) Clear() {
	var empty []uint
	s.words = empty
}

func (s *IntSet) Copy() *IntSet {
	var wordCopy []uint
	wordCopy = append(wordCopy, s.words...)
	copy := IntSet{words: wordCopy}
	return &copy
}

func (s *IntSet) AddAll(newWords ...int) {
	for _, word := range newWords {
		s.Add(word)
	}
}
