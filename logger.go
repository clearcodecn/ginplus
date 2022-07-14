// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginplus

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"time"
)

const requestIdKey = "x-request-id"

type logFormatter struct {
	*logrus.JSONFormatter
}

func (a *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	ctx := entry.Context
	// a. request-id
	requestId := ctx.Value("x-request-id")
	entry = entry.WithField("reqid", requestId)
	return a.JSONFormatter.Format(entry)
}

func DefaultLogger() *logrus.Entry {
	l := logrus.New()
	l.SetFormatter(&logFormatter{
		JSONFormatter: &logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339,
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
		},
	})
	e := logrus.NewEntry(l)
	return e
}

func LogContext(ctx context.Context) *logrus.Entry {
	return DefaultLogger().WithContext(ctx)
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig() HandlerFunc {
	sid := shortid.MustNew(1, "", uint64(time.Now().UnixNano()))
	return func(c *Context) {
		start := time.Now()
		if c.GetString(requestIdKey) == "" {
			c.Set(requestIdKey, sid.MustGenerate())
		}
		// Process request
		c.Next()
		dur := time.Since(start).Seconds()
		LogContext(c).
			WithFields(logrus.Fields{
				"cid":    Cid(c),
				"sec":    dur,
				"ip":     c.IP(),
				"ua":     c.Request.UserAgent(),
				"method": c.Request.Method,
				"url":    c.Request.URL.String(),
				"status": c.Writer.Status(),
				"size":   c.Writer.Size(),
			}).
			Infof("http_server_response")
	}
}
