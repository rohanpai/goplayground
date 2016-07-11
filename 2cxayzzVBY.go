package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
	tmpl  *template.Template
	i     = 0
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "test")
	tmpl.Execute(w, sess.Flashes())
	sess.Save(r, w)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "test")
	sess.AddFlash(fmt.Sprintf("Flash No %d", i))
	sess.Save(r, w)
	i++
	http.Redirect(w, r, "/", 302)
}

func main() {
	tmpl, _ = template.New("").Parse(`
<html>
	<body>
		<p>
			<a href="/test">test</a>
		</p>

		{{range .}}
			<p>{{.}}</p>
		{{end}}
	</body>
</html>`)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)
	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}
