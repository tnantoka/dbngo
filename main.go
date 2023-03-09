package main

import (
	"flag"
	"image/gif"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/tnantoka/dbngo/evaluator"
)

var input string
var outputPNG string
var outputGIF string
var scale int

func parseFlags() {
	flag.StringVar(&input, "i", "", "input file")
	flag.StringVar(&outputPNG, "p", "dbngo.png", "output png file")
	flag.StringVar(&outputGIF, "g", "", "output gif file")
	flag.IntVar(&scale, "s", 1, "scale")

	flag.Parse()

	if scale < 1 {
		log.Fatal("scale must be 1 or more")
	}
}

func main() {
	parseFlags()

	inputFile, err := os.Open(input)
	if err != nil {
		log.Fatalf("failed opening input file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPNG)
	if err != nil {
		log.Fatalf("failed creating output png file: %s", err)
	}
	defer outputFile.Close()

	e := evaluator.New()
	e.Scale = scale
	e.Directory = filepath.Dir(filepath.Clean(input))
	e.WithGIF = outputGIF != ""

	img := e.Eval(inputFile, input)

	if len(e.Errors) > 0 {
		log.Fatal(e.Errors)
	}

	if err := png.Encode(outputFile, img); err != nil {
		log.Fatalf("failed encoding image: %s", err)
	}

	if e.WithGIF {
		file, err := os.Create(outputGIF)
		if err != nil {
			log.Fatalf("failed creating output gif file: %s", err)
		}
		defer file.Close()
		if err := gif.EncodeAll(file, e.GIF); err != nil {
			log.Fatalf("failed encoding gif: %s", err)
		}
	}
}
