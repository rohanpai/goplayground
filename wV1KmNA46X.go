package main

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func main() {
	sessions.NewSession(store, "session-name")
	http.HandleFunc("/1", MyHandler1)
	http.HandleFunc("/2", MyHandler2)
	http.ListenAndServe(":8080", nil)
}

type Person struct {
	name string
}

type Monster struct {
	name string
}

func MyHandler1(w http.ResponseWriter, r *http.Request) {
	gob.Register(&Person{})
	gob.Register(&Monster{})

	session, _ := store.Get(r, "session-name")
	session.Values["person"] = &Person{"John Doe"}
	session.Values["monster"] = &Monster{"Sasquatch"}
	session.Save(r, w)
}

func MyHandler2(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	fmt.Println(session.Values["person"])
	fmt.Println(session.Values["monster"])
}
