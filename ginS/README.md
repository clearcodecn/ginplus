# Gin Default Server

This is API experiment for ginplus.

```go
package main

import (
	"github.com/clearcodecn/ginplus"
	"github.com/clearcodecn/ginplus/ginS"
)

func main() {
	ginS.GET("/", func(c *ginplus.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
