package main

import (
	"net/http"
	
	"appengine"

	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/bigquery/v2"
)

func init() {
	http.HandleFunc("/", hello)
}

func hello(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client, _ := serviceaccount.NewClient(c, bigquery.BigqueryScope)
	bqSvc, _ := bigquery.New(client)
	request := new(bigquery.TableDataInsertAllRequest)
	rows := make([]*bigquery.TableDataInsertAllRequestRows, 1)

	rows[0] = new(bigquery.TableDataInsertAllRequestRows)
	rows[0].Json = "<JSON ENTRY>"
	request.Rows = rows

	call := bqSvc.Tabledata.InsertAll("<PROJECTID>", "<dataset>", "<table>", request)
	_, err = call.Do()

	if err != nil {
		c.Errorf("error: %v", err)
	}
}
