package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

const BUFSIZE = 100
const RECORDSIZE = 255

type Record struct {
	Index byte
	Data  []byte
}

func (r *Record) ToBytes() []byte {
	out := make([]byte, 2, len(r.Data)+2)
	out[0] = r.Index
	out[1] = byte(len(r.Data))
	out = append(out, r.Data...)
	return out
}

func ReadToRecords(index byte, r io.ReadCloser, c chan Record, wg *sync.WaitGroup) {
	buf := make([]byte, RECORDSIZE)

	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Print(err)
			continue
		}
		if err == io.EOF {
			break
		}
		if n == 0 {
			continue
		}

		rec := Record{index, buf[:n]}
		c <- rec
	}
	wg.Done()
	r.Close()
}

func OutputRecords(c chan Record, w io.WriteCloser) {
	for r := range c {
		w.Write(r.ToBytes())
	}
	w.Close()
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Printf("Usage: mux <file1> <file2>...")
		return
	}
	if flag.NArg() > 255 {
		log.Printf("Too many files")
		return
	}

	c := make(chan Record, BUFSIZE)
	var wg sync.WaitGroup

	for index, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("Unable to open %s: %s", file, err)
			continue
		}
		wg.Add(1)
		go ReadToRecords(byte(index), file, c, &wg)
	}

	go OutputRecords(c, os.Stdout)

	wg.Wait()
	close(c)
}
