package converter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type Graphic struct {
	Name string
	HeadSize uint32
	BodySize uint32
	Offset int64
	Size uint32
	Buffer []byte
}

func Skater() ([]byte) {
	skater, err := os.ReadFile("./files/SKATER.P")
	if err != nil {
		if os.IsNotExist(err) {
			log.Panic("Skater file does not exist in files directory")
		}
	}
	idx := bytes.Index(skater, []byte("RW4ps3"))
	if idx != -1 {
		return skater
	} else {
		return nil
	}
}

func FindGraphicOffsets(skater []byte) ([]int) {
	var occurences []int
	data := skater
	term := []byte("RW4ps3")
	for x, d := bytes.Index(data, term), 0; x > -1; x, d = bytes.Index(data, term), d+x+1 {
		occurences = append(occurences, (x+d) - 1)
		data = data[x+1:]
	}
	return occurences
}

func ReadGraphics(skater []byte) ([]Graphic) {
	r := bytes.NewReader(skater)
	var graphics []Graphic

	for _, offset := range FindGraphicOffsets(skater) {
		var graphic Graphic
		position, _ := r.Seek(int64(offset), io.SeekStart)
		graphic.Offset = position

		r.Seek(0x44, io.SeekCurrent)
		graphic.HeadSize = ReadBEUint32(r)

		r.Seek(0x24, io.SeekCurrent)
		graphic.BodySize = ReadBEUint32(r)
		graphic.Size = graphic.HeadSize + graphic.BodySize

		r.Seek(0x14C, io.SeekCurrent)
		strBuffer := make([]byte, 8)
		br, _ := r.Read(strBuffer)
		graphic.Name = string(strBuffer[:br])

		buf := make([]byte, graphic.Size)
		r.Seek(graphic.Offset, io.SeekStart)
		rl, err := r.Read(buf)
		if err != nil {
			log.Panic()
		}
		graphic.Buffer = buf[:rl]
		graphics = append(graphics, graphic)
	}
	return graphics
}

func WriteGraphic(graphic Graphic) (error) {
	err := os.WriteFile(fmt.Sprintf("./files/%s.psg", graphic.Name), graphic.Buffer, 0666)
	if err != nil {
		return err
	}
	return nil
}