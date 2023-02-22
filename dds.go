package main

import (
	"bytes"
	"errors"
	"io"
)

type DDSFile struct {
	Size   int
	Width  uint32
	Height uint32
	Buffer []byte
}

func DDSParse(dds []byte) (DDSFile, error) {
	var file DDSFile
	if string(dds[:3]) == "DDS" && len(dds) > 128 {
		header := dds[:128] //We'll atleast assume there's a header currently there
		reader := bytes.NewReader(header)
		reader.Seek(0x0C, io.SeekStart)
		file.Height = ReadLEUint32(reader)
		file.Width = ReadLEUint32(reader)
		file.Size = len(dds)
		file.Buffer = dds
		return file, nil
	} else {
		return DDSFile{}, errors.New("file given does not meet DDS criteria")
	}
}
