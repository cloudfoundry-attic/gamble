Gamble (libyaml for GO)
===============================

Making the world better, one YAML parser at a time
--------------------------------------------------

Usage:

```go
package main

import (
  "fmt"
	"github.com/cloudfoundry/gamble"
)

var myYAML = `
---
some_key:
- some
- items
`

func main() {
	document, err := gamble.Parse(myYAML)
	if err != nil {
	  println("Could not parse yaml")
	  return
	}

  fmt.Printf("%#v\n", document)
  // prints out a map from string to interface
  // eg:
	//    map[string]interface{}{
	//    	"some_key": []interface{}{
	//    		"some",
	//    		"items",
	//    	},
	//    }
}

```
