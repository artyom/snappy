// Command snappy provides a basic utility to compress/decompress files/streams using snappy
// compression algorithm
package main

import (
	"flag"
	"io"
	"log"
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
	autoflags.Define(&p)
	flag.Parse()
	if p.Unpack {
		if err := decompress(p.From, p.To); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := compress(p.From, p.To); err != nil {
		log.Fatal(err)
	}
}

func compress(fromName, toName string) error {
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

func init() { log.SetFlags(0) }
