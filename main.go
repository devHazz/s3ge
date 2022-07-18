package main

import (
	"converter"
	//"fmt"
)

func main() {
	skater := converter.Skater()
	graphics := converter.ReadGraphics(skater)
	for _, graphic := range graphics {
		converter.WriteGraphic(graphic)
	}
}