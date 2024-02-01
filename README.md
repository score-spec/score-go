# score-go

Reference library for the parsing and loading SCORE files in Go.

This can be added to your project via:

```sh
$ go get -u github.com/score-spec/score-go
```

**NOTE**: if you project is still using the hand-written types, you will need to stay on `github.com/score-spec/score-go@v0.0.1`
and any important fixes to the schema may be back-ported to that branch.

## Parsing SCORE files

This library includes a few utility methods to parse source SCORE files.

```go
import (
    "io"
    "os"
    
    "github.com/score-spec/score-go/v1/loader"
    score "github.com/score-spec/score-go/v1/types"
)

func main() {
    var (
        err error
        src io.Reader
    )

    if src, err = os.Open("score.yaml"); err != nil {
        panic(err)
    }
	defer src.Close()
    
    var srcMap map[string]interface{}
    if err = loader.ParseYAML(&srcMap, src); err != nil {
        panic(err)
    }
    
    var spec score.WorkloadSpec
    if err = loader.MapSpec(&spec, srcMap); err != nil {
        panic(err)
    }

    // Do something with the spec
    // ...
}

```

## Upgrading the schema version

When the Score JSON schema is updated in https://github.com/score-spec/schema, this repo should be updated to match.

First update the subtree:

```
git subtree pull --prefix schema/files git@github.com:score-spec/schema.git main --squash
```

Then regenerate the defined types:

```
go generate -v ./...
```

And ensure the tests still pass:

```
go test -v ./...
```
