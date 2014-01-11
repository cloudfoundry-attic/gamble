Gamble (libyaml for GO)
===============================

Usage:

```go
package main

import (
	"github.com/cloudfoundry/gamble"
)

var myYAML = `
---
some_key:
- some
- items
`

func main() {
	document := gamble.Parse(myYAML)
	println(document == map[string]interface{}{
		"some_key": []interface{}{
			"some",
			"items",
		},
	})
}

```
