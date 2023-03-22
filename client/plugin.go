package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	pluginName       = "http-censor"
	ClientRegisterer = registerer(pluginName)
	licenseOK        = false
)

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", pluginName))
}

func (r registerer) registerClients(ctx context.Context, _ map[string]interface{}) (http.Handler, error) {
	logPrefix := fmt.Sprintf("[PLUGIN: %s]", pluginName)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		censor := false
		// we censor the gbp and we return the rupee :)
		if strings.Contains(r.URL.Path, "gbp") {
			r.URL.Path = strings.ReplaceAll(r.URL.Path, "gbp", "inr")
			censor = true
		}
		if strings.Contains(r.URL.RawQuery, "gbp") {
			r.URL.RawQuery = strings.ReplaceAll(r.URL.RawQuery, "gbp", "inr")
			censor = true
		}

		logger.Debug(fmt.Sprintf("%s - CENSORING %t : %s => %s", logPrefix,
			censor, r.URL.Path, r.URL.RawQuery))

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			logger.Error(logPrefix, err.Error())
		}

		for h, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		if resp.Body != nil {
			io.Copy(w, resp.Body)
			resp.Body.Close()
		}
	}), nil
}
