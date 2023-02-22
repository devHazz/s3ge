package main

import (
	"bytes"
	"fmt"
	"io"
)

var SkaterGraphicMagic = []byte{'\x89', '\x52', '\x57', '\x34', '\x70', '\x73', '\x33'}

type GraphicReader struct {
	reader *bytes.Reader
}

type Graphic struct {
	Id                       string
	Headsize, Bodysize, Size uint32
	Offset                   int64
}

func NewGraphic(r *bytes.Reader) *GraphicReader { return &GraphicReader{r} }

func (g GraphicReader) Read() (Graphic, error) {
	magic := make([]byte, 7)
	io.ReadFull(g.reader, magic)
	if !bytes.Equal(magic, SkaterGraphicMagic) {
		return Graphic{}, fmt.Errorf("invalid magic, got: %x", magic)
	}
	graphic := Graphic{}
	g.reader.Seek(0x3D, io.SeekCurrent)
	graphic.Headsize = ReadBEUint32(g.reader)

	g.reader.Seek(0x24, io.SeekCurrent)
	graphic.Bodysize = ReadBEUint32(g.reader)
	graphic.Size = graphic.Headsize + graphic.Bodysize

	g.reader.Seek(0x14C, io.SeekCurrent)
	id := make([]byte, 8)
	io.ReadFull(g.reader, id)
	graphic.Id = string(id)

	return graphic, nil
}

func (g GraphicReader) ReadMany() ([]Graphic, error) {

	return nil, nil
}
