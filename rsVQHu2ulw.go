package main

import (
	&#34;fmt&#34;
	&#34;net/http&#34;	
	&#34;appengine&#34;
	&#34;appengine/datastore&#34;
	&#34;appengine/memcache&#34;
)

func init() {
    http.HandleFunc(&#34;/&#34;, 					Index)
    http.HandleFunc(&#34;/without_flushall&#34;, 	QueryAfterPut_without_flushall)
    http.HandleFunc(&#34;/with_flushall&#34;, 		QueryAfterPut_with_flushall)
}

func Index(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/html; charset=utf-8&#34;)
	fmt.Fprintf(w, &#34;&lt;center&gt;&#34;)
	fmt.Fprintf(w, &#34;&lt;h1&gt;Count or query immediately after put&lt;/h1&gt;&#34;)
	fmt.Fprintf(w, &#34;&lt;table width=&#39;60%&#39;&gt;&#34;)
	fmt.Fprintf(w, &#34;&lt;tr&gt;&lt;td style=&#39;color:blue&#39;&gt;Ok     &lt;/td&gt;&lt;td&gt;  QueryAfterPut_with_flushall    &lt;/td&gt;&lt;td&gt;  &lt;a href=&#39;/with_flushall&#39;    target=&#39;_blank&#39;&gt;/with_flushall   &lt;/a&gt;&lt;/td&gt;&lt;/tr&gt;&#34;)
	fmt.Fprintf(w, &#34;&lt;tr&gt;&lt;td style=&#39;color:red&#39;&gt; Error  &lt;/td&gt;&lt;td&gt;  QueryAfterPut_without_flushall &lt;/td&gt;&lt;td&gt;  &lt;a href=&#39;/without_flushall&#39; target=&#39;_blank&#39;&gt;/without_flushall&lt;/a&gt;&lt;/td&gt;&lt;/tr&gt;&#34;)
	fmt.Fprintf(w, &#34;&lt;table&gt;&lt;/center&gt;&#34;)
}

func QueryAfterPut_without_flushall(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/html; charset=utf-8&#34;)
	type Demo struct {
		Key 			*datastore.Key 	`datastore:&#34;-&#34; json:&#34;-&#34;`
		Cat_byte		[]byte			
		Cat_string		string			
	}
	Key := datastore.NewKey(c, &#34;Demo&#34;,  &#34;StringID&#34;, 0, nil)
	
	e := &amp;Demo{}
	e.Key = Key
	e.Cat_byte 	 = []byte(&#34;1234&#34;)
	e.Cat_string = &#34;1234&#34;
	key, err := datastore.Put(c, e.Key, e)
	if err != nil {
		fmt.Fprintf(w, `put err: %s  &lt;br&gt;`,err)
		return
	}	
	fmt.Fprintf(w, `put success&lt;br&gt;`)

	count, err := datastore.NewQuery(&#34;Demo&#34;).KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  &lt;br&gt;`,count)
	
	err = datastore.Delete(c, key)
	if err != nil {
		fmt.Fprintf(w, `del err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `del success&lt;br&gt;`)

	count, err = datastore.NewQuery(&#34;Demo&#34;).KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  &lt;br&gt;`,count)
}	

func QueryAfterPut_with_flushall(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set(&#34;Content-Type&#34;, &#34;text/html; charset=utf-8&#34;)
	type Demo struct {
		Key 			*datastore.Key 	`datastore:&#34;-&#34; json:&#34;-&#34;`
		Cat_byte		[]byte			
		Cat_string		string			
	}
	Key := datastore.NewKey(c, &#34;Demo&#34;,  &#34;StringID&#34;, 0, nil)
	
	e := &amp;Demo{}
	e.Key = Key
	e.Cat_byte 	 = []byte(&#34;1234&#34;)
	e.Cat_string = &#34;1234&#34;
	key, err := datastore.Put(c, e.Key, e)
	if err != nil {
		fmt.Fprintf(w, `put err: %s  &lt;br&gt;`,err)
		return
	}	
	fmt.Fprintf(w, `put success&lt;br&gt;`)

	memcache.Flush(c)

	count, err := datastore.NewQuery(&#34;Demo&#34;).KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  &lt;br&gt;`,count)
	
	err = datastore.Delete(c, key)
	if err != nil {
		fmt.Fprintf(w, `del err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `del success&lt;br&gt;`)

	memcache.Flush(c)

	count, err = datastore.NewQuery(&#34;Demo&#34;).KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  &lt;br&gt;`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  &lt;br&gt;`,count)
}	
