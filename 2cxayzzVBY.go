package main

import (
	&#34;fmt&#34;
	&#34;html/template&#34;
	&#34;net/http&#34;

	&#34;github.com/gorilla/context&#34;
	&#34;github.com/gorilla/securecookie&#34;
	&#34;github.com/gorilla/sessions&#34;
)

var (
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
	tmpl  *template.Template
	i     = 0
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, &#34;test&#34;)
	tmpl.Execute(w, sess.Flashes())
	sess.Save(r, w)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, &#34;test&#34;)
	sess.AddFlash(fmt.Sprintf(&#34;Flash No %d&#34;, i))
	sess.Save(r, w)
	i&#43;&#43;
	http.Redirect(w, r, &#34;/&#34;, 302)
}

func main() {
	tmpl, _ = template.New(&#34;&#34;).Parse(`
&lt;html&gt;
	&lt;body&gt;
		&lt;p&gt;
			&lt;a href=&#34;/test&#34;&gt;test&lt;/a&gt;
		&lt;/p&gt;

		{{range .}}
			&lt;p&gt;{{.}}&lt;/p&gt;
		{{end}}
	&lt;/body&gt;
&lt;/html&gt;`)
	http.HandleFunc(&#34;/&#34;, indexHandler)
	http.HandleFunc(&#34;/test&#34;, testHandler)
	err := http.ListenAndServe(&#34;:8080&#34;, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}
