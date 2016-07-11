package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;net/http&#34;
	&#34;net/url&#34;
)

const (
	GET    = &#34;GET&#34;
	POST   = &#34;POST&#34;
	PUT    = &#34;PUT&#34;
	DELETE = &#34;DELETE&#34;
)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type ResourceBase struct{}

func (ResourceBase) Get(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, &#34;&#34;
}

func (ResourceBase) Post(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, &#34;&#34;
}

func (ResourceBase) Put(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, &#34;&#34;
}

func (ResourceBase) Delete(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, &#34;&#34;
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
	portString := fmt.Sprintf(&#34;:%d&#34;, port)
	http.ListenAndServe(portString, nil)
}

// TEST RESOURCE IMPLEMENTATION

type Test struct {
	// Default implementation of all Resource methods
	ResourceBase
}

// Override the Get method
func (t Test) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, &#34;YAY&#34;
}

func main() {
	var a Test
	AddResource(a, &#34;/&#34;)
	Start(4000)
}
