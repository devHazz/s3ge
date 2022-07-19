package importer

import (
	//"bytes"
	"converter"
	"log"
	"os"
)

func HandlePSG(index int, file string) (error) {
	newGraphic, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	s := converter.Skater()
	graphic := converter.ReadGraphics(s)[index]
	
	offset := converter.FindGraphicOffsets(s)[index]
	copy(s[offset + int(graphic.HeadSize):offset + int(graphic.Size)], newGraphic[graphic.HeadSize:graphic.Size])
	err = os.WriteFile("./files/SKATER.P", s, 0666)
	if err != nil {
		return err
	}
	return nil
}