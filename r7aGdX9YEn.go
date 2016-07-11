package main

import (
	&#34;encoding/base64&#34;
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;net/http&#34;
	&#34;net/http/httptest&#34;
)

func main() {
	recorder := httptest.NewRecorder()

	// Drop a cookie on the recorder.
	SetPreferencesCookie(recorder, &amp;Preferences{Colour: &#34;Blue&#34;})

	// Copy the Cookie over to a new Request
	request := &amp;http.Request{Header: http.Header{&#34;Cookie&#34;: recorder.HeaderMap[&#34;Set-Cookie&#34;]}}

	// Extract the dropped cookie from the request.
	cookie, err := request.Cookie(&#34;test&#34;)
	if err != nil {
		panic(&#34;Expected Cookie named &#39;test&#39;&#34;)
	}

	prefs, err := Decode(cookie)
	if err != nil {
		panic(&#34;Failed to decode cookie: &#34; &#43; string(err.Error()))
	}
	fmt.Printf(&#34;&gt;&gt; Decoded: %&#43;v\n&#34;, prefs)
}

type Preferences struct {
	Colour string
}

func SetPreferencesCookie(w http.ResponseWriter, pref *Preferences) error {
	data, err := json.Marshal(pref)
	if err != nil {
		return err
	}
	http.SetCookie(w, &amp;http.Cookie{
		Name:  &#34;test&#34;,
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
	if err := json.Unmarshal(data, &amp;prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}
