package main

import (
	"github.com/alexflint/go-arg"
	"log"
)

var Args struct {
	Input   string `arg:"required"`
	Convert bool   `arg:"-c,--convert"`
	Import  bool   `arg:"-i,--import"`
}

func main() {
	p := arg.MustParse(&Args)
	if !Args.Import && !Args.Convert {
		p.Fail("you must provide a flag for either importer or converter")
	}

	if Args.Import {
		p.Fail("importer not properly implemented")
	}

	if Args.Convert {
		err := Convert(Args.Input)
		if err != nil {
			log.Panic(err)
		}
	}
}
