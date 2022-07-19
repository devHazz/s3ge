package main

import (
	"converter"
	"fmt"
)

func main() {
	skater := converter.Skater()
	graphics := converter.ReadGraphics(skater)
	for _, graphic := range graphics {
		path := fmt.Sprintf("./files/%s.psg", graphic.Name)
		fmt.Printf("Path: %s | Header Size: %d | Body Size: %d | Offset: %d\n", path, graphic.HeadSize, graphic.BodySize, graphic.Offset)
		converter.WriteGraphic(graphic)
	}
}