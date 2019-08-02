package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminHandler_Listen_Records(t *testing.T) {
	ah := AdminHandler{}
	c := make(chan *ServedInfo, 2)
	ah.Listen(c, 1)

	c <- &ServedInfo{
		Request:  &RequestInfo{
			Method:  "GET",
			Path:    "/foo",
		},
		Response: &ResponseInfo{
			StatusCode: 200,
		},
	}

	c <- &ServedInfo{
		Request:  &RequestInfo{
			Method:  "POST",
			Path:    "/fred",
		},
		Response: &ResponseInfo{
			StatusCode: 400,
		},
	}

	// We could be waiting for approx twice the interval
	start := time.Now()
	for len(ah.requests) < 2 && time.Since(start).Nanoseconds() < time.Second.Nanoseconds() * 20 {
		time.Sleep(2)
	}

	if len(ah.requests) < 2 || time.Since(start).Nanoseconds() > time.Second.Nanoseconds() * 20 {
		t.Error("Assert timed out.")
		return
	}

	assertEqual(t, "POST", ah.requests[0].Request.Method)
	assertEqual(t, "/fred", ah.requests[0].Request.Path)
	assertEqual(t, 400, ah.requests[0].Response.StatusCode)

	assertEqual(t, 2, len(ah.requests))
	assertEqual(t, "GET", ah.requests[1].Request.Method)
	assertEqual(t, "/foo", ah.requests[1].Request.Path)
	assertEqual(t, 200, ah.requests[1].Response.StatusCode)

	close(c)
}

func TestAdminHandler_GetRequests(t *testing.T) {
	req, _ := http.NewRequest("GET", "/request", nil)
	res := httptest.NewRecorder()
	ah := AdminHandler{
		requests: []*ServedInfo{{
			Request: &RequestInfo{
				Method: "FRED",
				Path:   "/foo",
			},
			Response: &ResponseInfo{
				StatusCode: 123,
			},
		}},
	}
	handler := http.HandlerFunc(ah.ServeHTTP)

	handler.ServeHTTP(res, req)

	responseBody, _ := ioutil.ReadAll(res.Body)
	var responseModel []*ServedInfo
	json.Unmarshal(responseBody, &responseModel)

	assertEqual(t, "FRED", responseModel[0].Request.Method)
	assertEqual(t, "/foo", responseModel[0].Request.Path)
	assertEqual(t, 123, responseModel[0].Response.StatusCode)
}

func TestAdminHandler_DeleteRequests(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/request", nil)
	res := httptest.NewRecorder()
	ah := AdminHandler{
		requests: []*ServedInfo{{
			Request: &RequestInfo{
				Method: "FRED",
				Path:   "/foo",
			},
			Response: &ResponseInfo{
				StatusCode: 123,
			},
		}},
	}
	handler := http.HandlerFunc(ah.ServeHTTP)

	handler.ServeHTTP(res, req)

	assertEqual(t, 0, len(ah.requests))
}