set -eu

if [[ $CF_PAGES == '1' ]]; then
  go install golang.org/x/tools/cmd/goyacc@latest
  $GOPATH/bin/goyacc -o ../parser/parser.go ../parser/parser.go.y
fi

rm -rf main.wasm examples wasm_exec.js examples.js

cp -r ../examples ./examples

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
GOOS=js GOARCH=wasm go build -o main.wasm main.go
echo "const examples = [$(ls ../examples | awk '{ print "{ \"name\": \"" $1 "\" }," }')" > examples.js
echo "$(ls ../examples/mitpress | awk '{ print "{ \"name\": \"mitpress/" $1 "\" }," }')];" >> examples.js
