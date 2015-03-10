package main

import (
	"encoding/binary"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"syscall"

	"github.com/LiptonB/mux"
	"github.com/golang/protobuf/proto"
)

const BUFSIZE = 10

func ReadToRecords(recordsize int64, index uint32, r io.ReadCloser,
	c chan *mux.Record, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, recordsize)
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Print(err)
			continue
		}
		if n == 0 {
			if err == io.EOF {
				//log.Printf("EOF reached")
				break
			} else {
				log.Printf("Non-EOF empty read")
				continue
			}
		}

		rec := &mux.Record{
			Index: proto.Uint32(index),
			Data:  buf[:n],
		}
		c <- rec
	}
	wg.Done()
	r.Close()
}

func OutputRecords(c chan *mux.Record, w io.WriteCloser) {
	for r := range c {
		//log.Printf("%d bytes for %d", len(r.Data), *r.Index)
		serialized, err := proto.Marshal(r)
		if err != nil {
			log.Printf("Unable to serialize record, skipping")
			continue
		}
		binary.Write(w, binary.LittleEndian, uint32(len(serialized)))
		w.Write(serialized)
	}
	w.Close()
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Printf("Usage: mux <file1> <file2>...")
		return
	}
	if flag.NArg() > math.MaxUint32 {
		log.Printf("Too many files")
		return
	}

	c := make(chan *mux.Record, BUFSIZE)
	var wg sync.WaitGroup

	var stat syscall.Stat_t
	err := syscall.Fstat(int(os.Stdout.Fd()), &stat)
	if err != nil {
		log.Print(err)
		return
	}
	recordsize := stat.Blksize

	for index, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("Unable to open %s: %s", file, err)
			continue
		}
		wg.Add(1)
		go ReadToRecords(recordsize, uint32(index), file, c, &wg)
	}

	go OutputRecords(c, os.Stdout)

	wg.Wait()
	close(c)
}
