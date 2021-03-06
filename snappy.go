// Command snappy provides a basic utility to compress/decompress files/streams using snappy
// compression algorithm
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/artyom/autoflags"
	"github.com/golang/snappy"
)

func main() {
	p := struct {
		From   string `flag:"f,input file"`
		To     string `flag:"o,output file"`
		Unpack bool   `flag:"d,decompress"`
	}{}
	autoflags.Parse(&p)
	var fn func(string, string) error = compress
	if p.Unpack {
		fn = decompress
	}
	if err := fn(p.From, p.To); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func compress(fromName, toName string) error {
	var dst io.WriteCloser = os.Stdout
	var src io.ReadCloser = os.Stdin
	var err error
	if fromName != "" {
		if src, err = os.Open(fromName); err != nil {
			return err
		}
		defer src.Close()
	}

	if toName != "" {
		if dst, err = os.Create(toName); err != nil {
			return err
		}
		defer dst.Close()
	}
	w := snappy.NewBufferedWriter(dst)
	if _, err = io.Copy(w, src); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return dst.Close()
}

func decompress(fromName, toName string) error {
	var dst io.WriteCloser = os.Stdout
	var src io.ReadCloser = os.Stdin
	var err error
	if fromName != "" {
		src, err = os.Open(fromName)
	}
	if err != nil {
		return err
	}
	defer src.Close()

	if toName != "" {
		dst, err = os.Create(toName)
	}
	defer dst.Close()
	if _, err = io.Copy(dst, snappy.NewReader(src)); err != nil {
		return err
	}
	return dst.Close()
}
