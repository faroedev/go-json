# github.com/pilcrowonpaper/go-json

A JSON parser and encoder.

## Example

### Values

```go
package main

import (
    "fmt"
    "github.com/pilcrowonpaper/go-json"
)

func main() {
	// Parse any JSON-encoded string.
	jsonValue, err := json.Parse(data)
	if err != nil {
		log.Fatalf("Invalid json: %s", err.Error())
	}

	switch definedJSONValue := jsonValue.(type) {
	case json.StringType:
		fmt.Printf("string value: %s\n", string(definedJSONValue)) // wraps a string
	case json.NumberType:
		fmt.Printf("number value: %d\n", definedJSONValue.Int()) // wraps a string
	case json.BooleanType:
		fmt.Printf("boolean value: %t\n", definedJSONValue) // wraps a bool
	case json.NullType:
		fmt.Printf("null value\n")
	case json.ObjectType:
		// see Objects example
	case json.ObjectType:
		// see Arrays example
	}

	// Encode any JSON value to string.
	encoded := json.Encode(jsonValue)
}
```

### Objects

```go
package main

import (
    "fmt"
    "github.com/pilcrowonpaper/go-json"
)

func main() {
    jsonObject, err := json.ParseObject(data)
    if err != nil {
    	log.Fatal("invalid json or not a json object")
    }

    name, err := jsonObject.GetString("name")
    if errors.Is(err, json.ErrObjectMemberNotFound) {
    	log.Fatal("member not found")
    }
    if err != nil {
    	log.Fatal("value not a string")
    }
}
```

### Arrays

```go
package main

import (
    "fmt"
    "github.com/pilcrowonpaper/go-json"
)

func main() {
	jsonObject, err := json.ParseObject(data)
    if err != nil {
    	log.Fatal("invalid json or not a json object")
    }

    jsonName, err := jsonObject.GetString(0)
    if errors.Is(err, json.ErrArrayIndexOutOfBounds) {
    	log.Fatal("out of bounds")
    }
    if err != nil {
    	log.Fatal("value not a string")
    }
}
```

### Builder

```go
package main

import (
    "fmt"
    "github.com/pilcrowonpaper/go-json"
)

func main() {
    jsonObjectBuilder := json.NewObjectBuilder()
    jsonObjectBuilder.Add("name", StringType("pilcrow"))
    s := jsonObjectBuilder.Done()
    fmt.Println(s)
}
```
