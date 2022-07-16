package ginplus

func flash(ctx *Context) func(s string) string {
	return func(key string) string {
		return ctx.Flash(key)
	}
}

func hasSession(ctx *Context) func(s string) bool {
	return func(s string) bool {
		_, ok := ctx.session.Values[s]
		return ok
	}
}
