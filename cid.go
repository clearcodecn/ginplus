package ginplus

import (
	"context"
	"net"
	"strings"
)

const (
	cidKey       = "gin-plus-cid"
	cidHeaderKey = "x-cid"
)

const (
	configKey = "gin-plus-config"
)

// ContextWithCid 从 nginx 透传下来.这边不需要用域名去做解析.
func ContextWithCid() HandlerFunc {
	return func(ctx *Context) {
		ctx.Set(cidHeaderKey, ctx.GetHeader(cidHeaderKey))
	}
}

func Cid(ctx *Context) string {
	val, ok := ctx.Get(cidHeaderKey)
	if ok {
		return val.(string)
	}
	return ""
}

type Config struct {
	Robots      string `yaml:"robots"`      // 机器人
	SitemapFile string `yaml:"sitemapFile"` // 站点地图
	AdsTxt      string `yaml:"adsTxt"`      // 广告文件
}

func (e *Engine) SetConfig(config map[string]Config) {
	e.configs = config
}

func ContextWithConfig() HandlerFunc {
	return func(ctx *Context) {
		cid := Cid(ctx)
		config := ctx.engine.configs[cid]
		ctx.Set(configKey, config)
	}
}

func GetConfig(ctx *Context) Config {
	val, ok := ctx.Get(configKey)
	if ok {
		return val.(Config)
	}
	return Config{}
}

// IP 获取客户端 IP
func (c *Context) IP() string {
	if ip := c.Request.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	return ip
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}
