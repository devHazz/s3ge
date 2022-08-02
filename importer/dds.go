package importer

import (
	"bytes"
	"converter"
	"errors"
	"io"
)

type DDSFile struct {
	Magic string
	Size int
	Width uint32
	Height uint32
	Buffer []byte
}

func DDSParse(dds []byte) (DDSFile, error) {
	var file DDSFile
	file.Magic = string(dds[:3])
	if string(file.Magic) == "DDS" && len(dds) > 128 {
		//Check for the header magic and make sure there's atleast a header
		header := dds[:128] //We'll atleast assume there's a header currently there
		reader := bytes.NewReader(header)
		reader.Seek(0x0C, io.SeekStart)
		file.Height = converter.ReadLEUint32(reader)
		file.Width = converter.ReadLEUint32(reader)
		file.Size = len(dds)
		file.Buffer = dds
		return file, nil
	} else {
		return DDSFile{}, errors.New("file given does not meet DDS criteria")
	}
}