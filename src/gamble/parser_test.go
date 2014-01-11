package gamble

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var stringDocument = `
the_string
`

func TestParsesSingleString(t *testing.T) {
	node, err := Parse(stringDocument)
	assert.NoError(t, err)
	assert.Equal(t, node, "the_string")
}

var sequenceDocument = `
- foo
- bar
- baz
`

func TestParseSequences(t *testing.T) {
	node, err := Parse(sequenceDocument)
	assert.NoError(t, err)
	assert.Equal(t, node, []interface{} {
		"foo",
		"bar",
		"baz",
	})
}

var mapDocument = `
key1: value1
key2: value2
`

func TestParsesSingleMap(t *testing.T) {
	node, err := Parse(mapDocument)
	assert.NoError(t, err)
	assert.Equal(t, node, map[string]interface{} {
		"key1": "value1",
		"key2": "value2",
	})
}

var nestedMapDocument = `
---
globals:
- taco
- burrito
- kimchi
collections:
- name: oceans
  locals:
    foo: bar
    bar: baz
  sequences:
  - one
  - two
  - three
- name: seas
  age: 55
`

func TestParseNestedMap(t *testing.T) {
	node, err := Parse(nestedMapDocument)
	assert.NoError(t, err)
	assert.Equal(t, node, map[string]interface{}{
		"globals": []interface{} {
			"taco",
			"burrito",
			"kimchi",
		},
		"collections":[]interface{} {
			map[string]interface{} {
				"name": "oceans",
				"locals": map[string]interface{} {
					"foo": "bar",
					"bar": "baz",
				},
				"sequences": []interface{}{
					"one",
					"two",
					"three",
				},
			},
			map[string]interface{} {
				"name": "seas",
				"age": "55",
			},
		},
	})
}

var invalidDocument = `
---
-
	-
		-
`

func TestReturnsErrorForInvalidDocument(t *testing.T) {
	_, err := Parse(invalidDocument)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Error parsing YAML.")
}
