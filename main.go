package main

import (
	"converter"
	"fmt"
	"flag"
)

func main() {
	importbool := flag.Bool("i", false, "importer boolean")
	convertbool := flag.Bool("c", false, "converter boolean")

	if *importbool && !*convertbool {
		
	}

	skater := converter.Skater()
	graphics := converter.ReadGraphics(skater)
	for _, graphic := range graphics {
		path := fmt.Sprintf("./files/%s.psg", graphic.Name)
		fmt.Printf("Path: %s | Header Size: %d | Body Size: %d | Offset: %d\n", path, graphic.HeadSize, graphic.BodySize, graphic.Offset)
		converter.WriteGraphic(graphic)
	}
}