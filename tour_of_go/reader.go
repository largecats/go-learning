package main

import "golang.org/x/tour/reader"

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.

func (r MyReader) Read(b []byte) (int, error) {
	n := 0
	for ; n < len(b); n++ {
		b[n] = 'A'
	}
	return n, nil
}

func main() {
	reader.Validate(MyReader{})
}
