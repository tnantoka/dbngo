package renderer

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestRenderer(t *testing.T) {
	tests := []struct {
		input    [][]color.Color
		expected string
	}{
		{
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
			"white.png",
		},
		{
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}},
				{color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}},
			},
			"checked.png",
		},
	}

	for i, tt := range tests {
		img := Render(1, tt.input)

		debugFile, _ := os.Create(fmt.Sprintf("../tmp/test_%d.png", i))
		defer debugFile.Close()
		png.Encode(debugFile, img)

		file, err := os.Open("../testdata/" + tt.expected)
		if err != nil {
			t.Fatalf("failed to open file: %s", err)
		}
		defer file.Close()

		expected, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatalf("failed to read file: %s", err)
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			t.Fatalf("failed to encode image: %s", err)
		}
		actual := buf.Bytes()

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}
