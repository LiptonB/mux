package mux

import (
	"bufio"
	"errors"
	"io"
)

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

func RecordFromReader(r *bufio.Reader) (*Record, error) {
	index, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	// log.Printf("Found index: %d", index)

	length, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	// log.Printf("Found length: %d", length)

	var left int = int(length)
	buf := make([]byte, length)
	bufleft := buf

	for left > 0 {
		n, err := r.Read(bufleft)
		if err == io.EOF {
			return nil, errors.New("Unexpected EOF")
		} else if err != nil {
			return nil, err
		}
		left -= n
		bufleft = bufleft[n:]
		// log.Printf("Read %d bytes into buffer", n)
	}

	rec := &Record{index, buf}
	return rec, nil
}
