package ginplus

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/gorilla/sessions"
	"strings"
)

var (
	SessionKey  string
	SessionName = "laravel-session"
)

func Session() HandlerFunc {
	if SessionKey == "" {
		panic("session key is not set")
	}
	var (
		sessionStore = sessions.NewCookieStore([]byte(SessionKey))
	)
	return func(ctx *Context) {
		sess, _ := sessionStore.Get(ctx.Request, SessionName)
		ctx.session = sess
		ctx.beforeRender = func() {
			_ = sess.Save(ctx.Request, ctx.Writer)
		}
		ctx.Next()
	}
}

func (c *Context) AddFlash(key string, val string) {
	c.session.AddFlash(val, key)
}

func (c *Context) Flash(key string) string {
	val := c.session.Flashes(key)
	if len(val) > 0 {
		return val[0].(string)
	}
	return ""
}

func (c *Context) Session(key string, val string) {
	c.session.Values[key] = val
}

func (c *Context) SessionObject(key string, val interface{}) {
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(val)
	c.session.Values[key] = buf.String()
}

func (c *Context) SessionGet(key string) string {
	val, ok := c.session.Values[key]
	if !ok {
		return ""
	}
	return val.(string)
}

var (
	SessionNotFound = errors.New("session not found")
)

func (c *Context) SessionGetObj(key string, out interface{}) error {
	val, ok := c.session.Values[key]
	if !ok {
		return SessionNotFound
	}
	return gob.NewDecoder(strings.NewReader(val.(string))).Decode(out)
}

func (c *Context) SessionDel(key string) {
	delete(c.session.Values, key)
}
