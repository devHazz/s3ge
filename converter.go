package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"time"
)

func FindGraphicOffsets(skater []byte) []int {
	var offsets []int
	data := skater
	for x, d := bytes.Index(data, SkaterGraphicMagic), 0; x > -1; x, d = bytes.Index(data, SkaterGraphicMagic), d+x+1 {
		offsets = append(offsets, (x+d)-1)
		data = data[x+1:]
	}
	return offsets
}

func Convert(file string) error {
	f, e := os.ReadFile(file)
	if e != nil {
		if os.IsNotExist(e) {
			return e
		}
	}
	reader := bytes.NewReader(f)
	curTime := strconv.Itoa(int(time.Now().UnixMilli()))
	e = os.Mkdir(curTime, 0666)
	if e != nil {
		return e
	}
	for _, o := range FindGraphicOffsets(f) {
		pos, e := reader.Seek(int64(o)+1, io.SeekStart)
		g, e := NewGraphic(reader).Read()
		if e != nil {
			return e
		}
		e = os.WriteFile(curTime+"/"+g.Id+".psg", f[pos:pos+int64(g.Size)], 0666)
		if e != nil {
			return e
		}
	}
	return nil
}
