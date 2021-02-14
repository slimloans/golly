package golly

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// WebContext specific context for web
// this will allow us not to pass down Context
type WebContext struct {
	Context
	request *http.Request
	writer  http.ResponseWriter

	rendered bool

	urlParams map[string]string
}

// NewWebContext returns a new web context
func NewWebContext(a Application, r *http.Request, w http.ResponseWriter) WebContext {
	ctx := NewContext(r.Context())
	ctx.SetDB(a.DB)

	ctx.SetLogger(a.Logger.WithFields(webLogParams(ctx, r)))

	return WebContext{
		urlParams: map[string]string{},
		Context:   ctx,
		request:   r,
		writer:    w,
	}
}

func webLogParams(ctx Context, r *http.Request) log.Fields {
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logFields["http.proto"] = r.Proto
	logFields["http.request_id"] = ctx.ContextID()

	logFields["http.method"] = r.Method
	logFields["http.useragent"] = r.UserAgent()
	logFields["http.url"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	logFields["url"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	logFields["http.url_details.path"] = r.URL.Path
	logFields["http.url_details.host"] = r.Host
	logFields["http.url_details.queryString"] = r.URL.RawQuery
	logFields["http.url_details.schema"] = scheme

	return logFields
}

func (wctx WebContext) Request() *http.Request {
	return wctx.request
}

func (wctx WebContext) Writer() http.ResponseWriter {
	return wctx.writer
}

func (wctx *WebContext) setURLParams(params map[string]string) {
	wctx.urlParams = params
}

// URLParam returns a URL parameter
func (wctx *WebContext) URLParam(key string) string {
	return wctx.urlParams[key]
}

// AddHeader adds a reaponse header
func (wctx *WebContext) AddHeader(key, value string) {
	wctx.Writer().Header().Add(key, value)
}

// RenderStatus renders out a status
func (wctx *WebContext) RenderStatus(status int) {
	wctx.rendered = true
	wctx.Writer().WriteHeader(status)
}