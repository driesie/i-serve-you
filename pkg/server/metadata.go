// Copyright (c) 2019 Oliver Wyman Digital
// Use of this source code is governed by a MIT Licence
// license that can be found in the LICENSE file.

package server

type ServedInfo struct {
	Request  *RequestInfo
	Response *ResponseInfo
}

type ResponseInfo struct {
	StatusCode int
	Headers    map[string][]string
	Payload    []byte
}

type RequestInfo struct {
	Method  string
	Path    string
	Headers map[string][]string
	// converting payload to string for readability
	Payload string
}
