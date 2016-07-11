package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type ResourceBase struct{}

func (ResourceBase) Get(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Post(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Put(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Delete(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		request.ParseForm()
		method := request.Method
		values := request.Form

		switch method {
		case GET:
			code, data = resource.Get(values)
		case POST:
			code, data = resource.Post(values)
		case PUT:
			code, data = resource.Put(values)
		case DELETE:
			code, data = resource.Delete(values)
		default:
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func AddResource(resource Resource, path string) {
	http.HandleFunc(path, requestHandler(resource))
}

func Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}

// TEST RESOURCE IMPLEMENTATION

type Test struct {
	// Default implementation of all Resource methods
	ResourceBase
}

// Override the Get method
func (t Test) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, "YAY"
}

func main() {
	var a Test
	AddResource(a, "/")
	Start(4000)
}
