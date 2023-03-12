package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"strings"
	"syscall/js"

	"github.com/tnantoka/dbngo/evaluator"
)

func main() {
	js.Global().Set("generate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from", r)
			}
		}()

		input := args[0].String()
		return generate(input)
	}))

	c := make(chan struct{})
	<-c
}

func generate(input string) string {
	e := evaluator.New()

	img := e.Eval(strings.NewReader(input), "input")

	if len(e.Errors) > 0 {
		return strings.Join(e.Errors, "\n")
	}

	buf := &bytes.Buffer{}
	err := png.Encode(buf, img)
	if err != nil {
		return err.Error()
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURL := fmt.Sprintf("data:image/png;base64,%s", encoded)
	return dataURL
}
