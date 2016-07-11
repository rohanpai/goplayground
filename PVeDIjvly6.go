/*
	Package is to support getting an IP address from a client&#39;s browser and returing IP Address
	It provides functionality comparable to http://jsonip.appspot.com
*/

package getipaddress

import (
	&#34;encoding/json&#34;
	&#34;net/http&#34;
)

func init() {
	http.HandleFunc(&#34;/GetClientIPAddress&#34;, handlerGetClientIPAddress)
}

func handlerGetClientIPAddress(w http.ResponseWriter, r *http.Request) {
	//Verify that request has an origin handler
	if r.Header.Get(&#34;Origin&#34;) == &#34;&#34; {
		http.Error(w, &#34;Cross domain requests require Origin header&#34;, http.StatusBadRequest)
		return
	}

	//Verify that the request method is a GET
	if r.Method != &#34;GET&#34; {
		http.Error(w, &#34;Cross domain request only supports GET&#34;, http.StatusBadRequest)
		return
	}

	//Set response headers
	w.Header().Add(&#34;Access-Control-Allow-Origin&#34;, &#34;*&#34;)
	w.Header().Add(&#34;Access-Control-Allow-Methods&#34;, &#34;GET&#34;)
	w.Header().Add(&#34;Content-Type&#34;, &#34;application/json&#34;)

	//Get IP address from request and set it in return data
	retData := &amp;getClientIPAddressReturn{IP: r.RemoteAddr}
	encodedRetData, err := json.Marshal(retData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Return the response to the user
	w.Write(encodedRetData)
}

type getClientIPAddressReturn struct {
	IP string
}

//Run on dev command: