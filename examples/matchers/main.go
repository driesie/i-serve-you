// Copyright (c) 2019 Oliver Wyman Digital
// Use of this source code is governed by a MIT Licence
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"github.com/lshift/i-serve-you/pkg/server"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// use defaults
	config := server.NewConfig()

	server.Start(config,

		// match URL return payload
		func(r *http.Request) (b bool, i func(w http.ResponseWriter)) {
			if r.URL.Path == "/foo" {
				return true, func(writer http.ResponseWriter) {
					writer.WriteHeader(http.StatusCreated)
					writer.Write([]byte(`
{
	"payload": "Hello World"
}
`))
				}
			}

			return false, nil
		},

		// match method and part URL return error
		func(request *http.Request) (b bool, i func(http.ResponseWriter)) {
			if request.Method == "PATCH" && strings.Contains(request.URL.Path, "bar") {
				return true, func (res  http.ResponseWriter) {
					res.WriteHeader(http.StatusBadRequest)
				}
			}
			return false, nil
		},

		// inspect the body
		func(r *http.Request) (b bool, i func(w http.ResponseWriter)) {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			// create new reader so it can be read again
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			var j interface{}
			err := json.Unmarshal(bodyBytes, &j)
			if nil != err {
				return false, nil
			}
			m := j.(map[string]interface{})

			if m["fred"] == "hello" {
				return true, func(writer http.ResponseWriter) {
					writer.WriteHeader(http.StatusInternalServerError)
				}
			}

			return false, nil
		},
	)
}
