package main

import (
	"github.com/clearcodecn/ginplus"
	"time"
)

func main() {
	g := ginplus.New()
	g.Use(ginplus.ContextWithCid())
	g.Use(ginplus.Logger())
	g.GET("/", func(ctx *ginplus.Context) {
		time.Sleep(300 * time.Millisecond)
		ginplus.LogContext(ctx).Infof("testingaaa")
		ctx.JSONP(200, ginplus.H{
			"200": "aaa",
		})
	})
	g.Run(":3123")
}
