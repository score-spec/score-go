# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

.PHONY: .ALWAYS
.ALWAYS:

## Update score schema
.PHONY: update-schema
update-schema:
	rm -fv schema/files/score-v1b1.json.modified
	git subtree pull --prefix schema/files git@github.com:score-spec/schema.git main --squash

## Generate types
.PHONY: generate
generate:
# Unfortunately struct generators don't know how to handle mixed properties and additional properties so we have to strip these out before we generate the structs.
# We still validate with the original specification though.
	jq 'walk(if type == "object" and .type == "object" and .additionalProperties == true and (.properties | type) == "object" then (del(.required) | del(.properties)) else . end)' schema/files/score-v1b1.json > schema/files/score-v1b1.json.modified
	go generate -v ./...