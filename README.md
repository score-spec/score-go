[![ci](https://github.com/score-spec/score-k8s/actions/workflows/ci.yaml/badge.svg)](https://github.com/score-spec/score-k8s/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/score-spec/score-go.svg)](https://pkg.go.dev/github.com/score-spec/score-go)

# score-go

Reference library containing common types and functions for building Score implementations in Go.

This can be added to your project via:

```sh
go get -u github.com/score-spec/score-go@latest
```

**NOTE**: if you project is still using the hand-written types, you will need to stay on `github.com/score-spec/score-go@v0.0.1`
and any important fixes to the schema may be back-ported to that branch.

## Packages

- `github.com/score-spec/score-go/schema` - Go constant with the json schema, and methods for validating a json or yaml structure against the schema.
- `github.com/score-spec/score-go/types` - Go types for Score workloads, services, and resources generated from the json schema.
- `github.com/score-spec/score-go/loader` - Go functions for loading the validated json or yaml structure into a workload struct. 
- `github.com/score-spec/score-go/framework`  - Common types and functions for Score implementations.

## Parsing SCORE files

This library includes a few utility methods to parse source SCORE files.

```go
import (
    "os"

    "gopkg.in/yaml.v3"
	
    scoreloader "github.com/score-spec/score-go/loader"
    scoreschema "github.com/score-spec/score-go/schema"
    scoretypes "github.com/score-spec/score-go/types"
)

func main() {
    src, err := os.Open("score.yaml")
    if err != nil {
        panic(err)
    }
    defer src.Close()

    var srcMap map[string]interface{}
    if err := yaml.NewDecoder(src).Decode(&srcMap); err != nil {
        panic(err)
    }
    
    if err := scoreschema.Validate(srcMap); err != nil {
        panic(err)
    }

    var spec scoretypes.Workload
    if err := scoreloader.MapSpec(&spec, srcMap); err != nil {
        panic(err)
    }
    
    if err := scoreloader.Normalize(&spec, "."); err != nil {
        panic(err)
    }

    // Do something with the spec
    // ...
}
```

## Building a Score implementation

[score-compose](https://github.com/score-spec/score-compose) is the reference Score implementation written in Go and using this library. If you'd like to write a custom Score implementation, use the functions in this library and the `score-compose` implementation as a Guide.

## Upgrading the schema version

When the Score JSON schema is updated in <https://github.com/score-spec/spec>, this repo should be updated to match.

First copy the new `score-v1b1.json` and `samples/` files from the spec repo, into [schema/files](schema/files).

Then regenerate the defined types:

```sh
make generate
```

And ensure the tests still pass:

```sh
go test -v ./...
```
