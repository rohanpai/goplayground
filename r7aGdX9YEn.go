package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	recorder := httptest.NewRecorder()

	// Drop a cookie on the recorder.
	SetPreferencesCookie(recorder, &Preferences{Colour: "Blue"})

	// Copy the Cookie over to a new Request
	request := &http.Request{Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	cookie, err := request.Cookie("test")
	if err != nil {
		panic("Expected Cookie named 'test'")
	}

	prefs, err := Decode(cookie)
	if err != nil {
		panic("Failed to decode cookie: " + string(err.Error()))
	}
	fmt.Printf(">> Decoded: %+v\n", prefs)
}

type Preferences struct {
	Colour string
}

func SetPreferencesCookie(w http.ResponseWriter, pref *Preferences) error {
	data, err := json.Marshal(pref)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "test",
		Value: base64.StdEncoding.EncodeToString(data),
	})
	return nil
}

func Decode(cookie *http.Cookie) (*Preferences, error) {
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	var prefs *Preferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}
