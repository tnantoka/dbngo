package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/tnantoka/dbngo/evaluator"
	"github.com/tnantoka/dbngo/renderer"
)

var input string
var output string
var scale int

func parseFlags() {
	flag.StringVar(&input, "i", "", "input file")
	flag.StringVar(&output, "o", "dbngo.png", "output file")
	flag.IntVar(&scale, "s", 1, "scale")

	flag.Parse()

	if scale < 1 {
		log.Fatal("scale must be 1 or more")
	}
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func main() {
	parseFlags()

	inputFile, err := os.Open(input)
	if err != nil {
		log.Fatalf("failed opening input file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(output)
	if err != nil {
		log.Fatalf("failed creating output file: %s", err)
	}
	defer outputFile.Close()

	e := evaluator.New()
	pixels := e.Eval(inputFile)

	if len(e.Errors) > 0 {
		log.Fatal(e.Errors)
	}

	img := renderer.Render(scale, pixels)
	if err := png.Encode(outputFile, img); err != nil {
		log.Fatalf("failed writing image: %s", err)
	}
}
