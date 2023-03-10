# dbngo

A tiny [Design By Numbers](https://dbn.media.mit.edu/) clone written in Go.  
Generate PNG adn GIF from `.dbn` files.

## Commands

- [x] Paper
- [x] Pen
- [x] Line
- [x] Set
- [x] Set (Dot)
- [x] Set (Copy)
- [x] Repeat
- [x] Same/Notsame
- [x] Smaller/Notsmaller
- [x] Command
- [x] Load
- [x] Number
- ~~Mouse~~
- ~~Forever~~
- ~~Key~~
- ~~Net~~
- ~~Time~~

## Built-in libraries

- Letters

## Diffs

Command | DBN | dbngo
--- | --- | ---
Load | `Load lib.dbn` | `Load "lib.dbn"`

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

## Acknowledgements

- [Design By Numbers](https://dbn.media.mit.edu/)
- [Writing An Interpreter In Go](https://interpreterbook.com/)

## Author

[tnantoka](https://twitter.com/tnantoka)
