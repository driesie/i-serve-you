// Copyright (c) 2019 Oliver Wyman Digital
// Use of this source code is governed by a MIT Licence
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Port      int
	AdminPort int
}

func NewConfig() *Config {
	c := &Config{}
	// defaults
	c.Port = 8080
	c.AdminPort = 8081

	return c
}

func Start(config *Config, matchers ...func(request *http.Request) (bool, func(http.ResponseWriter))) {
	finish := make(chan bool)
	requestChannel := make(chan *ServedInfo, 100)

	adminHandler := &AdminHandler{
		[]*ServedInfo{},
	}
	adminHandler.Listen(requestChannel)
	go func() {
		fmt.Printf("Starting server on port %v.\n", config.Port)
		http.ListenAndServe(fmt.Sprintf(":%d", config.Port), &IServeYouHandler{requestChannel, matchers})
	}()

	go func() {
		fmt.Printf("Starting admin server on port %v.\n", config.AdminPort)
		http.ListenAndServe(fmt.Sprintf(":%d", config.AdminPort), adminHandler)
	}()

	<-finish
}

type IServeYouHandler struct {
	requests chan *ServedInfo
	matchers []func(request *http.Request) (bool, func(http.ResponseWriter))
}

type ResponseProxy struct {
	proxy *ResponseInfo
	real  http.ResponseWriter
}

func (p ResponseProxy) Header() http.Header {
	return p.real.Header()
}

func (p ResponseProxy) Write(d []byte) (int, error) {
	p.proxy.Payload = d
	return p.real.Write(d)
}

func (p ResponseProxy) WriteHeader(statusCode int) {
	p.proxy.StatusCode = statusCode
	p.real.WriteHeader(statusCode)
}

func (p ResponseProxy) CopyHeaders() {
	p.proxy.Headers = p.real.Header()
}

func (h IServeYouHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	matched := false
	response := ResponseProxy{&ResponseInfo{}, res}
	for _, m := range h.matchers {
		if match, responser := m(req); match {
			fmt.Printf("Matching %s\n", req.URL.Path)
			responser(response)
			matched = true
			break
		}
	}

	if !matched {
		fmt.Sprintf("No matcher for %s, default return success\n", req.URL.Path)
		response.WriteHeader(http.StatusOK)
	}

	data, _ := ioutil.ReadAll(req.Body)
	response.CopyHeaders()
	info := &ServedInfo{
		Request: &RequestInfo{
			Method:  req.Method,
			Path:    req.URL.Path,
			Headers: req.Header,
			Payload: string(data),
		},
		Response: response.proxy,
	}
	h.requests <- info
}
