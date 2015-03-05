package mux

import (
	"bufio"
	"errors"
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

	length, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errors.New("Not enough bytes read")
	}

	rec := &Record{index, buf}
	return rec, nil
}
