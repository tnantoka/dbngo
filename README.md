# dbngo

A tiny [Design By Numbers](https://dbn.media.mit.edu/) clone written in Go.  
Generate PNG adn GIF from `.dbn` files.

## Example

```
// gradient.dbn

Repeat A 0 100 {
  Pen A
  Set Y (100 - A)
  Line 0 Y 100 Y
}
```

```
$ dbngo -i gradient.dbn -p gradient.png -g gradient.gif -s 2
```

![](docs/gradient.png)
![](docs/gradient.gif)

## Live demo with wasm

https://dbngo.tnantoka.com/

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
- [ ] Array

## Built-in libraries

- Letters

## Diffs

Command | DBN | dbngo
--- | --- | ---
Load | `Load lib.dbn` | `Load "lib.dbn"`

## Examples

- ~~amoebic~~
- ~~bandedclock~~
- [x] calculate
- [x] commands
- ~~dancinguy~~
- ~~dancingy~~
- [x] dots
- ~~grainsofrain~~
- ~~graymachine~~
- ~~headsortails~~
- ~~intersecting~~
- [x] line
- ~~looping~~
- ~~meeber~~
- ~~merging~~
- ~~nervousguy~~
- [ ] nesting1
- [ ] nesting2
- ~~painting~~
- [x] paper
- ~~paramecium~~
- ~~plaid~~
- ~~probing~~
- ~~quantitative~~
- [x] questions
- ~~raininglines~~
- ~~reactive~~
- [x] repeating
- [x] rocket2
- ~~rockettime~~
- ~~thehunt~~
- ~~time1~~
- ~~time2~~
- ~~tuftball~~
- [x] variable

## Development

```
$ go install golang.org/x/tools/cmd/goyacc@latest

$ goyacc -o parser/parser.go parser/parser.go.y

$ go fmt ./...

$ go test $(go list ./... | grep -v /wasm) -coverprofile=cover_broken.out && \
  cat cover_broken.out | grep -v yaccpar | grep -v .y > cover.out && \
  go tool cover -html=cover.out -o coverage.html

$ go run main.go -i testdata/hello.dbn -p tmp/dbngo.png -g tmp/dbngo.gif -s 2
```

## Acknowledgements

- [Design By Numbers](https://dbn.media.mit.edu/)
- [Writing An Interpreter In Go](https://interpreterbook.com/)

## Author

[tnantoka](https://twitter.com/tnantoka)
