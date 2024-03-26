// Copyright 2024 Humanitec
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

//go:generate go run github.com/atombender/go-jsonschema@v0.15.0 -v --schema-output=https://score.dev/schemas/score=types.gen.go --schema-package=https://score.dev/schemas/score=types --schema-root-type=https://score.dev/schemas/score=Workload ../schema/files/score-v1b1.json.modified
