package main

import (
	"errors"
	"fmt"
	_ "github.com/gin-gonic/gin"
	"strings"
)

func main() {}

var (
	pluginName         = "content-caser"
	ModifierRegisterer = registerer(pluginName)

	errUnkownType                = errors.New("unknown response type")
	errRegistererNotFound        = fmt.Errorf("%s plugin disabled: config not found", pluginName)
	logger                Logger = noopLogger{}
	logPrefix                    = fmt.Sprintf("[PLUGIN: %s]", pluginName)
	keySeparator                 = "."
)

type registerer string

func (r registerer) RegisterModifiers(f func(
	name string,
	modifierFactory func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r), r.responseModifierFactory, false, true)
}

func (registerer) responseModifierFactory(cfg map[string]interface{}) func(interface{}) (interface{}, error) {
	config, err := NewCaseModifierConfig(cfg)
	if err != nil {
		return nil
	}

	fnCase := strings.ToUpper
	if !config.UseUpper {
		fnCase = strings.ToLower
	}

	return func(input interface{}) (interface{}, error) {
		resp, ok := input.(ResponseWrapper)
		if !ok {
			return nil, errUnkownType
		}
		headers := resp.Headers()
		if headers == nil {
			headers = make(map[string][]string)
		}
		data := resp.Data()
		if data == nil || len(data) == 0 {
			headers["X-Content-Case"] = []string{"false"}
			return responseWrapper{
				data: data,
				// TODO: check that we can remove this...
				isComplete: resp.IsComplete(),
				io:         resp.Io(), // TODO: check what is this for ..
				headers:    headers,
				statusCode: resp.StatusCode(),
			}, nil
		}
		headers["X-Content-Case"] = []string{"true"}
		newData := recase(fnCase, data)
		return responseWrapper{
			data:       newData,
			isComplete: resp.IsComplete(),
			io:         resp.Io(),
			headers:    headers,
			statusCode: resp.StatusCode(),
		}, nil
	}
}

func recaseValue(fnCase func(string) string, data interface{}) interface{} {
	if strV, ok := data.(string); ok {
		return fnCase(strV)
	}

	if mV, ok := data.(map[string]interface{}); ok {
		rmV := make(map[string]interface{}, len(mV))
		for k, v := range mV {
			rmV[k] = recaseValue(fnCase, v)
		}
		return rmV
	}

	if aV, ok := data.([]interface{}); ok {
		raV := make([]interface{}, 0, len(aV))
		for _, i := range aV {
			raV = append(raV, recaseValue(fnCase, i))
		}
		return raV
	}
	return data
}

func recase(fnCase func(string) string, data map[string]interface{}) map[string]interface{} {
	if data == nil || len(data) == 0 {
		return data
	}
	vv := recaseValue(fnCase, data)
	if mv, ok := vv.(map[string]interface{}); ok {
		return mv
	}
	return data
}

func (registerer) RegisterLogger(in interface{}) {
	if l, ok := in.(Logger); ok {
		logger = l
		logger.Debug(logPrefix, "Logger loaded")
	}
}
