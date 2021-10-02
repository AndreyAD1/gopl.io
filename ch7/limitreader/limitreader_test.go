package main

import (
	"io"
	"strings"
	"testing"
)

func TestLimitReader_Positive (t *testing.T) {
	initialString := "some test string"
	stringReader := strings.NewReader(initialString)
	limitedReader := LimitReader(stringReader, 100)
	targetSlice := make([]byte, 5)
	n, err := limitedReader.Read(targetSlice)
	if err != nil {
		t.Fatalf("Invalid error: expect: nil, got: %v", err)
	}
	if n != 5 {
		t.Fatalf("Invalid read byte number: expect: %d, got: %d", 5, n)
	}
	if limitedReader.readLimit != 95 {
		t.Errorf("Invalid reader limit: expect: %d, got: %d", 95, limitedReader.readLimit)
	}
	if string(targetSlice) != initialString[:5] {
		t.Errorf("Invalid target slice: expect: %v, got: %v", initialString[:5], string(targetSlice))
	}
	targetSlice2 := make([]byte, 100)
	n, err = limitedReader.Read(targetSlice2)
	if err != io.EOF {
		t.Fatalf("Invalid error: expect: io.EOF, got: %v", err)
	}
	if n != len(initialString) - 5 {
		t.Fatalf("Invalid read byte number: expect: %d, got: %d", 5, n)
	}
	if limitedReader.readLimit != int64(95 - n) {
		t.Errorf("Invalid reader limit: expect: %d, got: %d", 95 - n, limitedReader.readLimit)
	}
	if string(targetSlice2[:n]) != initialString[5:] {
		t.Errorf("Invalid target slice: expect: %v, got: %v", targetSlice2[:n], initialString[5:])
	}
	targetSlice3 := make([]byte, 10)
	n, err = limitedReader.Read(targetSlice3)
	if err != io.EOF {
		t.Fatalf("Invalid error: expect: io.EOF, got: %v", err)
	}
	if n != 0 {
		t.Fatalf("Invalid read byte number: expect: %d, got: %d", 0, n)
	}
}

func TestLimitReader_ReachLimit(t *testing.T) {
	initialString := "some test string"
	stringReader := strings.NewReader(initialString)
	limitedReader := LimitReader(stringReader, 5)
	targetSlice := make([]byte, 100)
	n, err := limitedReader.Read(targetSlice)
	if err != io.EOF {
		t.Fatalf("Invalid error: expect: io.EOF, got: %v", err)
	}
	if n != 5 {
		t.Fatalf("Invalid read byte number: expect: %d, got: %d", 5, n)
	}
	if limitedReader.readLimit != 0 {
		t.Errorf("Invalid reader limit: expect: %d, got: %d", 0, limitedReader.readLimit)
	}
	targetSlice2 := make([]byte, 5)
	n, err = limitedReader.Read(targetSlice2)
	if err != io.EOF {
		t.Fatalf("Invalid error: expect: io.EOF, got: %v", err)
	}
	if n != 0 {
		t.Fatalf("Invalid read byte number: expect: %d, got: %d", 0, n)
	}
}