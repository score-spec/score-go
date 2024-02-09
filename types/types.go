package types

//go:generate go run github.com/atombender/go-jsonschema@v0.15.0 -v --schema-output=https://score.dev/schemas/score=types.gen.go --schema-package=https://score.dev/schemas/score=types --schema-root-type=https://score.dev/schemas/score=Workload ../schema/files/score-v1b1.json.modified
