package main

import "github.com/clearcodecn/ginplus"

func main() {
	g := ginplus.New()
	g.HTMLRender = ginplus.NewTemplateManager("examples/templates", ".gohtml", true)
	g.GET("/", func(ctx *ginplus.Context) {
		ctx.HTML(200, "a/index", nil)
	})
	g.Run(":1111")
}
