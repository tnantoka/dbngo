# dbngo

A tiny [Design By Numbers](https://dbn.media.mit.edu/) clone written in Go.

## Development

```
$ go install golang.org/x/tools/cmd/goyacc@latest

$ goyacc -o parser/parser.go parser/parser.go.y

$ go fmt ./...

$ go test ./... -coverprofile=cover_broken.out && \
  cat cover_broken.out | grep -v yaccpar | grep -v .y > cover.out && \
  go tool cover -html=cover.out -o coverage.html

$ go run main.go -i testdata/hello.dbn -p tmp/dbngo.png -g tmp/dbngo.gif -s 2
```
