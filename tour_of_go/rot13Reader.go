package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot13 *rot13Reader) Read(b []byte) (int, error) {
	n, err := rot13.r.Read(b)
	if err != nil {
		return n, err
	}
	for i := 0; i < n; i++ {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] = ((b[i]+13)-'a')%26 + 'a'
		} else if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] = ((b[i]+13)-'A')%26 + 'A'
		}
	}
	return n, nil
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
