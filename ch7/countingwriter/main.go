package main

import (
	"bytes"
	"fmt"
	"io"
)

type writerWrapper struct {
	wrappedWriter io.Writer
	byteCounter *int64
}

func (w writerWrapper) Write(p []byte) (n int, err error) {
	byteNumber, err := w.wrappedWriter.Write(p)
	*w.byteCounter += int64(byteNumber)
	return byteNumber, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var byteArray int64
	newWriter := writerWrapper{w, &byteArray}
	return newWriter, &byteArray
}

func main() {
	var byteBuffer []byte
	writer := bytes.NewBuffer(byteBuffer)
	writerWrapper, writerCounter := CountingWriter(writer)
	fmt.Println(*writerCounter)
	fmt.Fprintf(writerWrapper, "text")
	fmt.Println(*writerCounter)
	fmt.Fprintf(writerWrapper, "more text")
	fmt.Println(*writerCounter)
}