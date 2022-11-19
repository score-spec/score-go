# score-go
Reference library for the parsing and loading SCORE files

## Parsing SCORE files

This library includes a few utility methods to parse source SCORE files.

```go
import (
    "io"
    "os"
    
    "github.com/score-spec/score-go/loader"
    score "github.com/score-spec/score-go/types"
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
