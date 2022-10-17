# score-go
Reference library for the parsing and loading SCORE files

## Parsing SCORE files

This library includes a few utility methods to parse source SCORE files.

```go
import(
    "errors"
    "io"

    loader "github.com/score-spec/score-go/loader"
    score "github.com/score-spec/score-go/types"
)

var err error
var src []byte

if src, err = ioutil.ReadFile("score.yaml"); err != nil {
    panic(err)
}

var srcMap map[string]interface{}
if err = loader.ParseYAML(src, &srcMap); err != nil {
    panic(err)
}

var spec score.WorkloadSpec
if err = loader.ParseYAML(src, &srcMap); err != nil {
    panic(err)
}

// Do something with the spec
// ...

```
