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
	copy(s[offset:], newGraphic)

	copy(s[offset + 0x1BC:offset + 0x1BC + 8], []byte(graphic.Name)) //Change new graphic id to existing id
	err = os.WriteFile("./files/SKATER.P", s, 0666)
	if err != nil {
		return err
	}
	return nil
}