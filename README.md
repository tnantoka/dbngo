# dbngo

A tiny [Design By Numbers](https://dbn.media.mit.edu/) clone written in Go.

## Development

```
$ go install golang.org/x/tools/cmd/goyacc@latest

$ goyacc -o parser/parser.go parser/parser.go.y

$ go fmt ./...

$ go test ./... -coverprofile=cover.out && \
  go tool cover -html=cover.out -o coverage.html

$ go run main.go -i testdata/hello.dbn -o tmp/dbngo.png
```
