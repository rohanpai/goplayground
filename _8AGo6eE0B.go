// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package guestbook

import (
	&#34;io&#34;
	&#34;net/http&#34;
	&#34;text/template&#34;
	&#34;time&#34;

	&#34;appengine&#34;
	&#34;appengine/datastore&#34;
	&#34;appengine/memcache&#34;
	&#34;appengine/user&#34;
)

type Greeting struct {
	Seq     int64
	Author  string
	Content string
	Date    time.Time
}

const keyGreet = &#34;key_greet&#34;

func getSeq(ctx appengine.Context) (uint64, error) {
	ctx.Debugf(&#34;%s&#34;, &#34;getSeq&#34;)
	v, err := memcache.IncrementExisting(ctx, keyGreet, 1)

	if err == memcache.ErrCacheMiss {
		ctx.Debugf(&#34;%s&#34;, &#34;GetSeqNext ErrCacheMiss&#34;)

		i := uint64(0) //TODO get from dstore
		return memcache.Increment(ctx, keyGreet, 1, i)
	}

	if err != nil {
		ctx.Debugf(&#34;Memcache error - not cache miss:  %s&#34;, err.Error())
		return 0, err
	}

	ctx.Debugf(&#34;GetSeqNext received val: %d&#34;, v)
	return v, nil
}

func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/plain; charset=utf-8&#34;)
	io.WriteString(w, &#34;Not Found&#34;)
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/plain; charset=utf-8&#34;)
	io.WriteString(w, &#34;Internal Server Error&#34;)
	c.Errorf(&#34;%v&#34;, err)
}

var mainPage = template.Must(template.New(&#34;guestbook&#34;).Parse(
	`&lt;html&gt;&lt;body&gt;
{{range .}}
{{with .Author}}&lt;b&gt;{{.|html}}&lt;/b&gt;{{else}}An anonymous person{{end}}
on &lt;em&gt;{{.Date.Format &#34;3:04pm, Mon 2 Jan&#34;}}&lt;/em&gt;
wrote &lt;blockquote&gt;{{.Content|html}}&lt;/blockquote&gt;
{{end}}
&lt;form action=&#34;/sign&#34; method=&#34;post&#34;&gt;
&lt;div&gt;&lt;textarea name=&#34;content&#34; rows=&#34;3&#34; cols=&#34;60&#34;&gt;&lt;/textarea&gt;&lt;/div&gt;
&lt;div&gt;&lt;input type=&#34;submit&#34; value=&#34;Sign Guestbook&#34;&gt;&lt;/div&gt;
&lt;/form&gt;&lt;/body&gt;&lt;/html&gt;
`))

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != &#34;GET&#34; || r.URL.Path != &#34;/&#34; {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery(&#34;Greeting&#34;).Order(&#34;-Date&#34;).Limit(10)
	var gg []*Greeting
	_, err := q.GetAll(c, &amp;gg)
	if err != nil {
		serveError(c, w, err)
		return
	}
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/html; charset=utf-8&#34;)
	if err := mainPage.Execute(w, gg); err != nil {
		c.Errorf(&#34;%v&#34;, err)
	}
}

func handleSign(w http.ResponseWriter, r *http.Request) {
	if r.Method != &#34;POST&#34; {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		serveError(c, w, err)
		return
	}

	seq, err := getSeq(c)
	if err != nil {
		serveError(c, w, err)
		return
	}

	g := &amp;Greeting{
		Seq:     int64(seq),
		Content: r.FormValue(&#34;content&#34;),
		Date:    time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, &#34;Greeting&#34;, nil), g); err != nil {
		serveError(c, w, err)
		return
	}
	http.Redirect(w, r, &#34;/&#34;, http.StatusFound)
}

func init() {
	http.HandleFunc(&#34;/&#34;, handleMainPage)
	http.HandleFunc(&#34;/sign&#34;, handleSign)
}
