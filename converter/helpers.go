package converter

import (
	"encoding/binary"
	"io"
	"bytes"
)

func GetPosition(file bytes.Reader) int64 {
	offset, _ := file.Seek(0, io.SeekCurrent)
	return offset
}

func ReadBEUint32(r io.Reader) uint32 {
	buf := make([]byte, 4)
	io.ReadFull(r, buf)
	return binary.BigEndian.Uint32(buf)
}

func ReadLEUint32(r io.Reader) uint32 {
	buf := make([]byte, 4)
	io.ReadFull(r, buf)
	return binary.LittleEndian.Uint32(buf)
}

func ReadBEUint16(r io.Reader) uint16 {
	buf := make([]byte, 2)
	io.ReadFull(r, buf)
	return binary.BigEndian.Uint16(buf)
}

func ReadLEUint16(r io.Reader) uint16 {
	buf := make([]byte, 2)
	io.ReadFull(r, buf)
	return binary.LittleEndian.Uint16(buf)
}
