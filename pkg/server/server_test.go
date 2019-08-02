package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIServeYouHandler_ServeHTTP_No_Matches_Return_Ok_And_Capture(t *testing.T) {
	testCases := []struct {
		method  string
		path    string
		headers map[string][]string
		payload []byte
	}{
		{"GET", "/foo", map[string][]string{"X-Test": {"test header"}}, nil},
		{"POST", "/fred", map[string][]string{"X-Test": {"another test header"}}, []byte("hello world")},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("[%s] %s", tc.method, tc.path), func(t *testing.T) {
			var body io.Reader
			if nil != tc.payload {
				body = bytes.NewReader(tc.payload)
			}
			req, _ := http.NewRequest(tc.method, tc.path, body)
			for k, vs := range tc.headers {
				for _, v := range vs {
					req.Header.Add(k, v)
				}
			}

			requestChannel := make(chan *ServedInfo, 2)
			isyh := IServeYouHandler{requests: requestChannel}
			isyh.matchers = append(isyh.matchers, func(request *http.Request) (bool, func(http.ResponseWriter)){
				return false,  nil
			})

			res := httptest.NewRecorder()
			handler := http.HandlerFunc(isyh.ServeHTTP)

			handler.ServeHTTP(res, req)
			capturedRequest := <-requestChannel
			close(requestChannel)

			assertEqual(t, 200, res.Result().StatusCode)

			assertEqual(t, tc.method, capturedRequest.Request.Method)
			assertEqual(t, tc.path, capturedRequest.Request.Path)
			for k, vs := range tc.headers {
				for i, v := range vs {
					assertEqual(t, v, capturedRequest.Request.Headers[k][i])
				}
			}

			if nil == tc.payload {
				assertEqual(t, "", capturedRequest.Request.Payload)
			} else {
				assertEqual(t, string(tc.payload), capturedRequest.Request.Payload)
			}
			assertEqual(t, 200, capturedRequest.Response.StatusCode)
		})
	}
}

func TestIServeYouHandler_ServeHTTP_Use_First_Matcher(t *testing.T) {
	req, _ := http.NewRequest("POST", "/foo", nil)

	matcher1 := func(request *http.Request) (bool, func(http.ResponseWriter)){
		return false,  nil
	}
	matcher2 := func(request *http.Request) (bool, func(http.ResponseWriter)){
		return true,  func (res  http.ResponseWriter) {
			res.WriteHeader(http.StatusBadRequest)
		}
	}

	requestChannel := make(chan *ServedInfo, 2)
	isyh := IServeYouHandler{requests: requestChannel}
	isyh.matchers = append(isyh.matchers, matcher1, matcher2)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(isyh.ServeHTTP)

	handler.ServeHTTP(res, req)
	<-requestChannel
	close(requestChannel)

	assertEqual(t, http.StatusBadRequest, res.Result().StatusCode)
}

type FakeResponseWriter struct {
	header      http.Header
	body        []byte
	statusCode int
}

func (fw FakeResponseWriter) Write(data []byte) (int, error) {
	fw.body = data
	return 0, nil
}

func (fw FakeResponseWriter) WriteHeader(statusCode int) {
	fw.statusCode = statusCode
}

func (fw FakeResponseWriter) Header() http.Header {
	return fw.header
}

func TestResponseProxy_CopyHeaders(t *testing.T) {
	fake := FakeResponseWriter{}
	fake.header = http.Header{"X-Test": []string{"test"}}
	proxy := ResponseProxy{&ResponseInfo{}, fake}

	proxy.CopyHeaders()

	assertEqual(t, "test", proxy.proxy.Headers["X-Test"][0])
}

func TestResponseProxy_Header(t *testing.T) {
	fake := FakeResponseWriter{}
	fake.header = http.Header{"X-Test": []string{"test"}}
	proxy := ResponseProxy{&ResponseInfo{}, fake}

	assertEqual(t, "test", proxy.Header()["X-Test"][0])
}

func TestResponseProxy_Write(t *testing.T) {
	fake := FakeResponseWriter{}
	proxy := ResponseProxy{&ResponseInfo{}, fake}

	proxy.Write([]byte("hello world"))

	assertEqual(t, "hello world", string(proxy.proxy.Payload))
}
