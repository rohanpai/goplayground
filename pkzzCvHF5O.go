package main

import (
	&#34;net/http&#34;
	
	&#34;appengine&#34;

	&#34;code.google.com/p/goauth2/appengine/serviceaccount&#34;
	&#34;code.google.com/p/google-api-go-client/bigquery/v2&#34;
)

func init() {
	http.HandleFunc(&#34;/&#34;, hello)
}

func hello(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client, _ := serviceaccount.NewClient(c, bigquery.BigqueryScope)
	bqSvc, _ := bigquery.New(client)
	request := new(bigquery.TableDataInsertAllRequest)
	rows := make([]*bigquery.TableDataInsertAllRequestRows, 1)

	rows[0] = new(bigquery.TableDataInsertAllRequestRows)
	rows[0].Json = &#34;&lt;JSON ENTRY&gt;&#34;
	request.Rows = rows

	call := bqSvc.Tabledata.InsertAll(&#34;&lt;PROJECTID&gt;&#34;, &#34;&lt;dataset&gt;&#34;, &#34;&lt;table&gt;&#34;, request)
	_, err = call.Do()

	if err != nil {
		c.Errorf(&#34;error: %v&#34;, err)
	}
}
