/*
	Package is to support getting an IP address from a client's browser and returing IP Address
	It provides functionality comparable to http://jsonip.appspot.com
*/

package getipaddress

import (
	"encoding/json"
	"net/http"
)

func init() {
	http.HandleFunc("/GetClientIPAddress", handlerGetClientIPAddress)
}

func handlerGetClientIPAddress(w http.ResponseWriter, r *http.Request) {
	//Verify that request has an origin handler
	if r.Header.Get("Origin") == "" {
		http.Error(w, "Cross domain requests require Origin header", http.StatusBadRequest)
		return
	}

	//Verify that the request method is a GET
	if r.Method != "GET" {
		http.Error(w, "Cross domain request only supports GET", http.StatusBadRequest)
		return
	}

	//Set response headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Content-Type", "application/json")

	//Get IP address from request and set it in return data
	retData := &getClientIPAddressReturn{IP: r.RemoteAddr}
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