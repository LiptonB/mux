package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"github.com/LiptonB/mux"
)

const BUFSIZE = 100

func ReadRecords(r io.ReadCloser, out []chan *mux.Record) {
	br := bufio.NewReader(r)

	for {
		rec, err := mux.RecordFromReader(br)
		if err != nil {
			break
		}
		out[rec.Index] <- rec
	}

	for _, c := range out {
		close(c)
	}
}

func WriteStream(c chan *mux.Record, w io.WriteCloser) {
	for rec := range c {
		w.Write(rec.Data)
	}
	w.Close()
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Printf("Usage: demux <file1> <file2>...")
		return
	}
	if flag.NArg() > 255 {
		log.Printf("Too many files")
		return
	}

	cs := make([]chan *mux.Record, flag.NArg())
	for i := range cs {
		cs[i] = make(chan *mux.Record, BUFSIZE)
	}

	for index, filename := range flag.Args() {
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Unable to open %s: %s", file, err)
			continue
		}
		go WriteStream(cs[index], file)
	}

	ReadRecords(os.Stdin, cs)
}
