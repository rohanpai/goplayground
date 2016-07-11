package main

import (
	"fmt"
	"net/http"	
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

func init() {
    http.HandleFunc("/", 					Index)
    http.HandleFunc("/without_flushall", 	QueryAfterPut_without_flushall)
    http.HandleFunc("/with_flushall", 		QueryAfterPut_with_flushall)
}

func Index(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<center>")
	fmt.Fprintf(w, "<h1>Count or query immediately after put</h1>")
	fmt.Fprintf(w, "<table width='60%'>")
	fmt.Fprintf(w, "<tr><td style='color:blue'>Ok     </td><td>  QueryAfterPut_with_flushall    </td><td>  <a href='/with_flushall'    target='_blank'>/with_flushall   </a></td></tr>")
	fmt.Fprintf(w, "<tr><td style='color:red'> Error  </td><td>  QueryAfterPut_without_flushall </td><td>  <a href='/without_flushall' target='_blank'>/without_flushall</a></td></tr>")
	fmt.Fprintf(w, "<table></center>")
}

func QueryAfterPut_without_flushall(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	type Demo struct {
		Key 			*datastore.Key 	`datastore:"-" json:"-"`
		Cat_byte		[]byte			
		Cat_string		string			
	}
	Key := datastore.NewKey(c, "Demo",  "StringID", 0, nil)
	
	e := &Demo{}
	e.Key = Key
	e.Cat_byte 	 = []byte("1234")
	e.Cat_string = "1234"
	key, err := datastore.Put(c, e.Key, e)
	if err != nil {
		fmt.Fprintf(w, `put err: %s  <br>`,err)
		return
	}	
	fmt.Fprintf(w, `put success<br>`)

	count, err := datastore.NewQuery("Demo").KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  <br>`,count)
	
	err = datastore.Delete(c, key)
	if err != nil {
		fmt.Fprintf(w, `del err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `del success<br>`)

	count, err = datastore.NewQuery("Demo").KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  <br>`,count)
}	

func QueryAfterPut_with_flushall(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	type Demo struct {
		Key 			*datastore.Key 	`datastore:"-" json:"-"`
		Cat_byte		[]byte			
		Cat_string		string			
	}
	Key := datastore.NewKey(c, "Demo",  "StringID", 0, nil)
	
	e := &Demo{}
	e.Key = Key
	e.Cat_byte 	 = []byte("1234")
	e.Cat_string = "1234"
	key, err := datastore.Put(c, e.Key, e)
	if err != nil {
		fmt.Fprintf(w, `put err: %s  <br>`,err)
		return
	}	
	fmt.Fprintf(w, `put success<br>`)

	memcache.Flush(c)

	count, err := datastore.NewQuery("Demo").KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  <br>`,count)
	
	err = datastore.Delete(c, key)
	if err != nil {
		fmt.Fprintf(w, `del err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `del success<br>`)

	memcache.Flush(c)

	count, err = datastore.NewQuery("Demo").KeysOnly().Count(c)
	if err != nil {
		fmt.Fprintf(w, `count err: %s  <br>`,err)
		return
	}
	fmt.Fprintf(w, `count: %d  <br>`,count)
}	
