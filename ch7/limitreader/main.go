package main

import (
	"fmt"
	"io"
	"strings"
)

type LimitedReader struct {
	reader io.Reader
	readLimit int64
}

func (l LimitedReader) Read(p []byte) (int, error) {
	if l.readLimit == 0 {
		return 0, io.EOF
	}

	if l.readLimit < int64(len(p)) {
		readByteNumber, _ := l.reader.Read(p[:l.readLimit])
		return readByteNumber, io.EOF
	} else {
		readByteNumber, err := l.reader.Read(p)
		l.readLimit -= int64(readByteNumber)
		return readByteNumber, err
	}
}

func LimitReader(r io.Reader, limit int64) io.Reader {
	return LimitedReader{r, limit}
}

func main() {
	stringReader := strings.NewReader("some test string")
	limitedReader := LimitReader(stringReader, 6)
	targetSlice := make([]byte, 20)
	limitedReader.Read(targetSlice)
	fmt.Println(targetSlice)
	fmt.Println(string(targetSlice))
}