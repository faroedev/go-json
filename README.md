# github.com/pilcrowonpaper/go-json

A JSON parser and encoder.

## Example

### Objects

```go
package main

import (
    "fmt"
    "github.com/pilcrowonpaper/go-json"
)

func main() {
    jsonObject, err := json.Parse(data)
    if err != nil {
        panic(err)
    }

    name, err := jsonObject.GetString("name")
    if err != nil {
        // Key doesn't exist or the value isn't a string.
        panic(err)
    }
    fmt.Println(name)

    jsonObject.SetString(name, "faroe")
    fmt.Println(jsonObject.String())
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
    jsonArray, err := json.ParseArray(data)
    if err != nil {
        panic(err)
    }

    name, err := jsonArray.GetString(0)
    if err != nil {
        // Item doesn't exist or the value isn't a string.
        panic(err)
    }
    fmt.Println(name)

    jsonArray.SetString(0, "faroe")
    fmt.Println(jsonArray.String())
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
    jsonObjectBuilder.AddString("name", "faroe")
    s := jsonObjectBuilder.Done()
    fmt.Println(s)
}
```
