package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var (
	pluginName        = "content-hash-etag"
	HandlerRegisterer = registerer(pluginName)
	errWrongConfig    = fmt.Errorf("%s: wrong config", pluginName)

	ifNoneMatch = "If-None-Match"
	eTag        = "ETag"
)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	logger.Info(fmt.Sprintf("[PLUGIN: %s]: Successfully registered", pluginName))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		knownEtag := requestHeader(r, ifNoneMatch)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		resp := rec.Result()
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var body []byte
		if resp.Body == nil {
			body = []byte("")
		} else {
			body = rec.Body.Bytes()
		}
		sum := sha256.Sum256(body)
		curETag := base64.StdEncoding.EncodeToString(sum[:])

		// There can be "weak" eTag (only taking into
		// account the boyd), or "strong" eTag taking
		// into account the headers: for this test we
		// only use the weak one:
		wHeaders := w.Header()
		copyHeaders(resp.Header, wHeaders)
		wHeaders.Set("ETag", string(curETag))
		if knownEtag == curETag {
			// n.ot sure if we should spit only this header..
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}), nil
}

func copyHeaders(src http.Header, dst http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func requestHeader(r *http.Request, key string) string {
	if values := r.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (r registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", pluginName))
}
