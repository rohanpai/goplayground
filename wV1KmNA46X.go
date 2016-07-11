package main

import (
	&#34;encoding/gob&#34;
	&#34;fmt&#34;
	&#34;net/http&#34;

	&#34;github.com/gorilla/securecookie&#34;
	&#34;github.com/gorilla/sessions&#34;
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func main() {
	sessions.NewSession(store, &#34;session-name&#34;)
	http.HandleFunc(&#34;/1&#34;, MyHandler1)
	http.HandleFunc(&#34;/2&#34;, MyHandler2)
	http.ListenAndServe(&#34;:8080&#34;, nil)
}

type Person struct {
	name string
}

type Monster struct {
	name string
}

func MyHandler1(w http.ResponseWriter, r *http.Request) {
	gob.Register(&amp;Person{})
	gob.Register(&amp;Monster{})

	session, _ := store.Get(r, &#34;session-name&#34;)
	session.Values[&#34;person&#34;] = &amp;Person{&#34;John Doe&#34;}
	session.Values[&#34;monster&#34;] = &amp;Monster{&#34;Sasquatch&#34;}
	session.Save(r, w)
}

func MyHandler2(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, &#34;session-name&#34;)
	fmt.Println(session.Values[&#34;person&#34;])
	fmt.Println(session.Values[&#34;monster&#34;])
}
