# score-go

Reference library for the parsing and loading SCORE files in Go.

This can be added to your project via:

```sh
go get -u github.com/score-spec/score-go@latest
```

**NOTE**: if you project is still using the hand-written types, you will need to stay on `github.com/score-spec/score-go@v0.0.1`
and any important fixes to the schema may be back-ported to that branch.

## Parsing SCORE files

This library includes a few utility methods to parse source SCORE files.

```go
import (
    "os"

    "github.com/score-spec/score-go/loader"
    "github.com/score-spec/score-go/schema"
    score "github.com/score-spec/score-go/types"
)

func main() {
    src, err := os.Open("score.yaml")
    if err != nil {
        panic(err)
    }
    defer src.Close()

    var srcMap map[string]interface{}
    if err := loader.ParseYAML(&srcMap, src); err != nil {
        panic(err)
    }
    
    if err := schema.Validate(srcMap); err != nil {
        panic(err)
    }

    var spec score.Workload
    if err := loader.MapSpec(&spec, srcMap); err != nil {
        panic(err)
    }
    
    if err := loader.Normalize(&spec, "."); err != nil {
        panic(err)
    }

    // Do something with the spec
    // ...
}
```

## Upgrading the schema version

When the Score JSON schema is updated in <https://github.com/score-spec/schema>, this repo should be updated to match.

First update the subtree:

```sh
make update-schema
```

Then regenerate the defined types:

```sh
make generate
```

And ensure the tests still pass:

```sh
go test -v ./...
```
