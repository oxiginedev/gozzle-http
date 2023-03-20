# Gozzle-HTTP

Minimalistic axios style HTTP client for golang

## Installation

```bash
go get github.com/adedaramola/gozzle-http
```

## Sample Usage

```go
package main

import (
    "fmt"

    "github.com/adedaramola/gozzle-http"
)

func main() {
    res, err := gozzle.Send(&gozzle.Config{
        URL: "https://jsonplaceholder.typicode.com/posts",
        Method: "POST",
        Body: gozzle.Map{
            "title": "Gozzle",
            "body": "Gozzle is the best golang http client",
            "userId": 1,
        },
    })

    if err != nil {
        fmt.Println(err)
    }

    fmt.Printf("%s\n",res.Data)
    fmt.Println(res.Status)
    fmt.Println(res.StatusText)
}
```
