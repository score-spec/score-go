# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

.PHONY: default
default:
	# Please provide a valid make target

## Generate types
.PHONY: generate
generate:
# Then the second modification is to remove the array version of the container files and volumes. Anything based on score-go will just handle the map type.
	jq '. as $$a | ."$$defs".container.properties.files |= $$a."$$defs".container.properties.files.oneOf[1] | ."$$defs".container.properties.volumes |= $$a."$$defs".container.properties.volumes.oneOf[1] | del(."$$defs".containerFile.properties.target) | del(."$$defs".containerVolume.properties.target)' schema/files/score-v1b1.json > schema/files/score-v1b1.json.for-validation
# Unfortunately struct generators don't know how to handle mixed properties and additional properties so we have to strip these out before we generate the structs.
# We still validate with the original specification though.
	jq 'walk(if type == "object" and .type == "object" and .additionalProperties == true and (.properties | type) == "object" then (del(.required) | del(.properties)) else . end)' schema/files/score-v1b1.json.for-validation > schema/files/score-v1b1.json.for-generation
	go generate -v ./...