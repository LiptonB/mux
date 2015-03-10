package mux

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

func RecordFromReader(r *bufio.Reader) (*Record, error) {
	var length uint32

	err := binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	//log.Printf("Found length: %d", length)

	buf := make([]byte, length)
	n, err := io.ReadFull(r, buf)
	if err == io.ErrUnexpectedEOF && n == 0 {
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	}

	rec := &Record{}
	err = proto.Unmarshal(buf, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}
