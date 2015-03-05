package mux

import "io"

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

func RecordFromReader(r io.ByteReader) (*Record, error) {
	rec := &Record{}

	return rec, nil
}
