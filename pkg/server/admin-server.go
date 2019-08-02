// Copyright (c) 2019 Oliver Wyman Digital
// Use of this source code is governed by a MIT Licence
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AdminHandler struct {
	requests []*ServedInfo
}

func (h *AdminHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// NOTE: not using a router as I'm trying to keep this free of 3rd party dependencies.
	//  this should be simple enough for the current use case
	switch req.URL.Path {
	case "/request":
		switch req.Method {
		case "GET":
			h.GetRequests(res, req)
		case "DELETE":
			h.DeleteRequests(res, req)
		}
	}
}

func (h *AdminHandler) GetRequests(res http.ResponseWriter, req *http.Request) {
	json, _ := json.Marshal(h.requests)
	res.Write(json)
}

func (h *AdminHandler) DeleteRequests(res http.ResponseWriter, req *http.Request) {
	h.requests = []*ServedInfo{}
	fmt.Fprintln(res, "Deleted")
}

func (h *AdminHandler) Listen(requestChannel <-chan *ServedInfo, interval time.Duration) {
	go func() {
		for {
			select {
			case request := <-requestChannel:
				if nil == request {
					return
				}

				fmt.Printf("Receiving %s\n", request.Request.Path)

				h.requests = append([]*ServedInfo{request}, h.requests...)
			default:
				time.Sleep(interval)
			}
		}
	}()
}