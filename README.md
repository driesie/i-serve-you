# I Serve You

I serve you is a test/fake/mock REST API server that will always serve you.
It exposes 2 endpoints, on different ports, one will accept any request, 
record this request and respond with a success response code unless configured otherwise.
The other returns all previous requests and responses for inspection and investigation and allows
for the log to be cleared.

It is a very simple Go program, with no 3rd party dependencies and can be used stand alone or
as part of another process, for example as part of an integration test suite.
The requests are recorded in memory and will disappear with the process. There is currently no
persistent backing.
It is meant to be used for testing and debugging and is purposely kept simple.

## Download/Install
```bash
go get TBC
```

Alternatively, git clone the repository.

Binary distribution will be provided in future.

## Usage
### Start Standalone server

Start the server with default config (ports 8080 and 8081):
```bash
./i-serve-you
```

Or specify ports, e.g.:
```bash
./i-serve-you -port 8085 -adminPort 8086
```
### Recording and inspecting

Once the server is running, you can send a request to the endpoint like so:
```bash
curl -XPATCH 'http://localhost:8085/something/here'
```
The server will respond successfully and this request can be seen  by calling the admin
endpoint like so:
```bash
curl -XGET 'http://localhost:8086/request'
```

And the server will respond like so:
```json
[
    {
        "Request": {
            "Method": "PATCH",
            "Path": "/something/here",
            "Headers": {
                "Accept": [
                    "*/*"
                ],
                "User-Agent": [
                    "curl/7.54.0"
                ]
            },
            "Payload": ""
        },
        "Response": {
            "StatusCode": 200,
            "Headers": {},
            "Payload": null
        }
    }
]
```

### Matchers
When the server is started, one or more matcher functions can be passed in. This will allow
you to specify specific responses to a given type of request

See [example](examples/matchers/example.go) of how to bootstrap such
server.

A matcher can be as simple or complicated as required, for example:

```go
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

		// respond based on request body content
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
```

### Embeded
See [example](examples/matchers/example.go)

## TODO
- CI and publish binaries for Mac, Windows and Linux

## Future enhancements
- A html/JavaScript interface for the admin endpoint would be nice
- Allowing matchers to be configured and persisted in config.
- Support for creating matchers using the admin UI above.
- Provide pre-defined matchers that support:
    - proxying requests
    - delay responses / variable response times
    - un-reliable responses (e.g. drop x% of requests)

Want to see anything else? Raise an Issue or Pull Request