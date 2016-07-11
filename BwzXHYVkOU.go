package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	//data type
	TypeString	= 1
	TypeBool	= 2
	TypeInt		= 3
	TypeInt32	= 4
	TypeInt16	= 5
	TypeFloat	= 6
	TypeDate	= 7
	TypeTime	= 8
	TypeGeo		= 9	//uses string format "lat,lon"

	DefQueryLimit		= 10
	EsHLFragment		= 1
	EsHLFragmentSize	= 100

	DateFormat	= "2006-01-02"
	EsDateFormat	= "YYYY-MM-dd"
)

type M map[string]interface{}
type nopCloser struct{ io.Reader }

func (nopCloser) Close() error	{ return nil }

var (
	Remap	bool	//re-map Riak buckets to reflect definition, (ES not supported, must create + reindex)

	Name				string
	Host, Port, ConnStr		string
	EsHost, EsPort, EsConnStr	string
	tables				map[string]*table_

	//ES highlighting
	EsHLTagPre, EsHLTagPost	[]string
)

type table_ struct {
	Name	string
	Key	string
	Cols	map[string]*tableCol_

	//elasticsearch
	Es	bool
}

type tableCol_ struct {
	Name	string
	Type	int
	Index	string
	Boost	float64
}

func init() {
	flag.BoolVar(&Remap, "dbremap", false, "Remap DB to reflect changes in definition")

	Host = "127.0.0.1"
	Port = "8098"
	EsHost = "127.0.0.1"
	EsPort = "9200"

	EsHLTagPre, EsHLTagPost = []string{"<strong>"}, []string{"</strong>"}
}

func Init(cfgFile string) error {
	//load config
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	m := M{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if v, ok := m["name"]; ok {
		Name = v.(string)
	}
	if v, ok := m["host"]; ok {
		Host = v.(string)
	}
	if v, ok := m["port"]; ok {
		Port = v.(string)
	}
	if v, ok := m["es_host"]; ok {
		EsHost = v.(string)
	}
	if v, ok := m["es_port"]; ok {
		EsPort = v.(string)
	}

	tables = make(map[string]*table_)
	if v, ok := m["tables"]; ok {
		tablesM := v.([]interface{})
		for _, v := range tablesM {
			t := v.(map[string]interface{})
			table := &table_{Name: t["name"].(string), Cols: make(map[string]*tableCol_)}
			if v, ok := t["key"]; ok {
				table.Key = v.(string)
			}
			if v, ok := t["es"]; ok {
				table.Es = v.(bool)
			}
			if v, ok := t["cols"]; ok {	//columns
				tc := v.([]interface{})
				for _, v := range tc {
					t := v.(map[string]interface{})
					col := &tableCol_{Name: t["name"].(string)}
					switch t["type"].(string) {
					case "string":
						col.Type = TypeString
					case "bool":
						col.Type = TypeBool
					case "int":
						col.Type = TypeInt
					case "int32":
						col.Type = TypeInt32
					case "int16":
						col.Type = TypeInt16
					case "float":
						col.Type = TypeFloat
					case "date":
						col.Type = TypeDate
					case "time":
						col.Type = TypeTime
					case "geo":
						col.Type = TypeGeo
					}

					if v, ok := t["index"]; ok {
						col.Index = v.(string)
					}
					if v, ok := t["boost"]; ok {
						col.Boost = v.(float64)
					}
					table.Cols[col.Name] = col
				}
			}
			tables[table.Name] = table
		}
	}

	//connections
	ConnStr = "http://" + Host + ":" + Port + "/"
	EsConnStr = "http://" + EsHost + ":" + EsPort + "/"

	//server ready?
	err = Ping()
	if err != nil {
		return err
	}

	//ensure tables match configs
	for k, _ := range tables {
		err = OpenTable(k)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	return nil
}

func Ping() error {	//to check if server online or not
	//riak -> GET /ping
	req, err := http.NewRequest("GET", ConnStr+"ping", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Riak not ready")
	}

	//ES -> GET /_cluster/health
	req, err = http.NewRequest("GET", EsConnStr, nil)
	if err != nil {
		return err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("ES request failed")
	}
	m, err := GetBodyJson(resp)
	if err != nil {
		return err
	}
	if !m["ok"].(bool) {
		return errors.New("ES not ready")
	}
	return nil
}

/* Table */
func OpenTable(table string) error {
	tbl, ok := tables[table]
	if !ok {
		return errors.New("Table not defined")
	}
	log.Println("Checking table " + table)

	//riak -> GET /buckets/bucket/props
	req, err := http.NewRequest("GET", ConnStr+"buckets/"+table+"/props", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if Remap || resp.StatusCode != 200 {	//create it
		//riak -> PUT /buckets/bucket/props
		req, err = http.NewRequest("PUT", ConnStr+"buckets/"+table+"/props", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		SetBodyString(req, `{"props":{"allow_mult":false, "last_write_wins":true}}`)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 204 {
			return errors.New("Riak fail to setup bucket")
		}
		log.Println("Riak bucket props for " + table + " created")
	}

	if tbl.Es {
		//ES -> GET /index/_mapping
		req, err = http.NewRequest("GET", EsConnStr+table+"/_mapping", nil)
		if err != nil {
			return err
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {	//not exists, create it
			err = esMapTable(tbl)
			if err != nil {
				return err
			}
			log.Println("ES mapping for " + table + " created")
		}
	}

	log.Println("Table " + table + " OK")
	return nil
}

func esMapTable(table *table_) error {
	mc := M{}
	for _, v := range table.Cols {
		c := M{"store": "yes"}	//store everything since we're using this as main DB
		switch v.Type {
		case TypeString:
			c["type"] = "string"
		case TypeBool:
			c["type"] = "boolean"
		case TypeFloat:
			c["type"] = "double"
		case TypeInt:
			c["type"] = "long"
		case TypeInt32:
			c["type"] = "integer"
		case TypeInt16:
			c["type"] = "short"
		case TypeDate:
			c["type"] = "date"
			c["format"] = EsDateFormat
		case TypeTime:
			c["type"] = "date"	//use ES default dateoptionaltime RFC3339
		case TypeGeo:
			c["type"] = "geo_point"
			c["lat_lon"] = true
		}
		if v.Index == "" || v.Index == "no" {
			c["index"] = "no"
		} else if v.Index == "yes" {
			//ES default to indexed
		} else if v.Index == "exact" {
			c["index"] = "not_analyzed"
		} else {	//fulltext, etc.
			c["index"] = "analyzed"
			c["analyzer"] = "ih_" + v.Index
		}

		if v.Boost > 0 {
			c["boost"] = v.Boost
		}
		mc[v.Name] = c
	}
	m := M{"doc": M{"ignore_conflicts": false, "properties": mc}}
	m = M{"mappings": m}

	req, err := http.NewRequest("PUT", EsConnStr+table.Name, nil)
	if err != nil {
		return err
	}
	err = SetBodyJson(req, m)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("ES mapping fail")
	}
	return nil
}

/* Data Object */
type Data_ struct {
	table	*table_
	m, mh	M
	Found	bool
}

func NewData(table string) *Data_ {
	return &Data_{table: tables[table], m: make(M)}
}
func (o *Data_) Clear() *Data_			{ o.m, o.mh, o.Found = make(M), nil, false; return o }
func (o *Data_) ColExists(col string) bool	{ _, ok := o.m[col]; return ok }
func (o *Data_) DeleteCol(col string) *Data_	{ delete(o.m, col); return o }
func (o *Data_) CloneTo(table string) *Data_ {
	d := NewData(table)
	for k, v := range o.m {
		d.m[k] = v
	}
	return d
}
func (o *Data_) GetMap() M	{ return o.m }	//use only for testing
func (o *Data_) GetJson() string {
	b, err := json.Marshal(o.m)
	if err == nil {
		return string(b)
	}
	return ""
}

func (o *Data_) Set(col string, v interface{}) *Data_ {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.String {	//add only if not empty
		if r.Len() > 0 {
			o.m[col] = v
		}
	}
	return o
}
func (o *Data_) SetStr(col string, v string) *Data_ {
	if len(v) > 0 {
		o.m[col] = v
	}
	return o
}
func (o *Data_) SetBool(col string, v bool) *Data_	{ o.m[col] = v; return o }
func (o *Data_) SetInt(col string, v int) *Data_	{ o.m[col] = float64(v); return o }
func (o *Data_) SetInt32(col string, v int32) *Data_	{ o.m[col] = float64(v); return o }
func (o *Data_) SetInt16(col string, v int16) *Data_	{ o.m[col] = float64(v); return o }
func (o *Data_) SetFloat(col string, v float64) *Data_	{ o.m[col] = v; return o }
func (o *Data_) SetTime(col string, v time.Time) *Data_ {
	if !v.IsZero() {
		o.m[col] = v
	}
	return o
}
func (o *Data_) SetTimeStr(col string, v string, layout string) *Data_ {
	if len(v) > 0 {
		t, err := time.Parse(layout, v)
		if err == nil {
			if !t.IsZero() {
				o.m[col] = t
			}	//add only if non-zero time
		}
	}
	return o
}
func (o *Data_) SetDate(col string, v string) *Data_	{ return o.SetTimeStr(col, v, DateFormat) }

func (o *Data_) Str(col string) string {
	v, ok := o.m[col]
	if ok {
		return v.(string)
	}
	return ""
}
func (o *Data_) Bool(col string) bool {
	v, ok := o.m[col]
	if ok {
		return v.(bool)
	}
	return false
}
func (o *Data_) Int(col string) int {
	v, ok := o.m[col]
	if ok {
		return int(v.(float64))
	}
	return 0
}
func (o *Data_) Int32(col string) int32 {
	v, ok := o.m[col]
	if ok {
		return int32(v.(float64))
	}
	return 0
}
func (o *Data_) Int16(col string) int16 {
	v, ok := o.m[col]
	if ok {
		return int16(v.(float64))
	}
	return 0
}
func (o *Data_) Float(col string) float64 {
	v, ok := o.m[col]
	if ok {
		return v.(float64)
	}
	return 0
}
func (o *Data_) Time(col string) time.Time {
	v, ok := o.m[col]
	if ok {
		return v.(time.Time)
	}
	return time.Time{}
}
func (o *Data_) Date(col string) string {
	v, ok := o.m[col]
	if ok {
		return v.(time.Time).Format(DateFormat)
	}
	return ""
}
func (o *Data_) Highlight(col string) string {
	if o.mh != nil {
		v, ok := o.mh[col]
		if ok {
			v := v.([]interface{})
			return v[0].(string)
		}
	}
	return ""
}

//Operations
func (o *Data_) Put() error	{ return o.put(false) }
func (o *Data_) CondPut() error	{ return o.put(true) }	//save only if not exists

func (o *Data_) put(cond bool) error {
	v, ok := o.m[o.table.Key]
	if !ok {
		return errors.New("Key is not defined")
	}

	if cond {	//only put if none
		err := o.Exists(v.(string))
		if err != nil {
			return err
		}
		if o.Found {
			return errors.New("Key already existed")
		}
	}

	id := url.QueryEscape(v.(string))
	j, err := json.Marshal(o.m)
	if err != nil {
		return err
	}

	//riak -> PUT /buckets/bucket/keys/key
	req, err := http.NewRequest("PUT", ConnStr+"buckets/"+o.table.Name+"/keys/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = nopCloser{bytes.NewBuffer(j)}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return errors.New("Riak store fail")
	}

	if o.table.Es {
		//ES -> PUT /tablename/doc/id
		req, err = http.NewRequest("PUT", EsConnStr+o.table.Name+"/doc/"+id, nil)
		if err != nil {
			return err
		}
		req.Body = nopCloser{bytes.NewBuffer(j)}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return errors.New("ES index failed")
		}
	}

	return nil
}

func (o *Data_) Exists(key string) error {	//Riak only
	//GET /buckets/bucket/keys/key
	req, err := http.NewRequest("GET", ConnStr+"buckets/"+o.table.Name+"/keys/"+url.QueryEscape(key), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		o.Found = true
	} else if resp.StatusCode != 404 {
		return errors.New("Riak fetch fail")
	}
	return nil
}

func (o *Data_) Get(key string, cols []string) error {	//Riak only, cols=nil means fetch all cols
	//GET /buckets/bucket/keys/key
	req, err := http.NewRequest("GET", ConnStr+"buckets/"+o.table.Name+"/keys/"+url.QueryEscape(key), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		m, err := GetBodyJson(resp)
		if err != nil {
			return err
		}
		o.m, o.Found = m, true
	} else if resp.StatusCode != 404 {
		return errors.New("Riak fetch fail")
	}
	return nil
}

func (o *Data_) Delete(key string) error {
	id := url.QueryEscape(key)

	//riak -> DELETE /buckets/bucket/keys/key
	req, err := http.NewRequest("DELETE", ConnStr+"buckets/"+o.table.Name+"/keys/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 && resp.StatusCode != 404 {
		return errors.New("Riak delete fail")
	}

	if o.table.Es {
		//ES -> DELETE /tablename/doc/id
		req, err = http.NewRequest("DELETE", EsConnStr+o.table.Name+"/doc/"+id, nil)
		if err != nil {
			return err
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			return errors.New("ES delete fail")
		}
	}

	return nil
}

/* Query Object */
type Query_ struct {
	table		*table_
	m, mf, ms, mh	M	//query, filter, sort, highlight
	cols		[]string
	offset, limit	int
}

func NewQuery(table string) *Query_ {
	return &Query_{table: tables[table], m: M{}, mf: M{}, limit: DefQueryLimit}
}
func (o *Query_) Clear() {
	o.m, o.mf, o.ms, o.mh, o.offset, o.limit = M{}, M{}, nil, nil, 0, DefQueryLimit
}
func (o *Query_) Cols(c ...string) *Query_	{ o.cols = c; return o }	//col to fetch
func (o *Query_) Offset(offset int) *Query_	{ o.offset = offset; return o }
func (o *Query_) Limit(limit int) *Query_	{ o.limit = limit; return o }
func (o *Query_) Sort(col string, dir string) *Query_ {
	if o.ms == nil {
		o.ms = M{}
	}
	o.ms[col] = dir
	return o
}
func (o *Query_) SortScore() *Query_	{ return o.Sort("_score", "desc") }

func (o *Query_) Highlight(col string, size int) *Query_ {
	if o.mh == nil {
		o.mh = M{}
	}
	if size == 0 {
		o.mh[col] = M{}
	} else {
		o.mh[col] = M{"fragment_size": size}
	}
	return o
}

func (o *Query_) Find() (*Result_, error) {	//ES only
	//construct query
	var m M
	if len(o.mf) > 0 {
		m = M{"filtered": M{"query": o.m, "filter": o.mf}}
	} else {
		m = o.m
	}
	m = M{"query": m, "from": o.offset, "size": o.limit}
	if len(o.cols) > 0 {
		m["fields"] = o.cols
	}
	if o.ms != nil {
		m["sort"] = o.ms
	}
	if o.mh != nil {
		m["highlight"] = M{"pre_tags":	EsHLTagPre, "post_tags":	EsHLTagPost,
			"number_of_fragments":	EsHLFragment, "fragment_size":	EsHLFragmentSize, "fields":	o.mh}
	}

	req, err := http.NewRequest("GET", EsConnStr+o.table.Name+"/doc/_search", nil)
	if err != nil {
		return nil, err
	}
	err = SetBodyJson(req, m)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("ES find fail")
	}
	m, err = GetBodyJson(resp)
	if err != nil {
		return nil, err
	}

	//compile results
	res := &Result_{}
	m = m["hits"].(map[string]interface{})
	res.Found = int(m["total"].(float64))
	hits := m["hits"].([]interface{})
	res.Count = len(hits)
	res.Data = make([]*Data_, res.Count)
	for i := 0; i < len(hits); i++ {
		m = hits[i].(map[string]interface{})
		d := &Data_{table: o.table, m: m["_source"].(map[string]interface{}), Found: true}
		if v, ok := m["highlight"]; ok {
			d.mh = v.(map[string]interface{})
		}
		res.Data[i] = d
	}

	return res, nil
}

//Collector
func (o *Query_) All() *Query_	{ o.m["match_all"] = M{}; return o }

func (o *Query_) Equal(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m["term"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m["term"] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Fulltext(col string, val interface{}, isAnd bool) *Query_ {
	var m M
	if v, ok := o.m["match"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m["match"] = m
	}
	if isAnd {
		m[col] = M{"query": val, "operator": "and"}
	} else {
		m[col] = val
	}
	return o
}

func (o *Query_) Prefix(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m["prefix"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m["prefix"] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Wildcard(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m["wildcard"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m["wildcard"] = m
	}
	m[col] = "*" + val.(string) + "*"
	return o
}

func (o *Query_) Range(col string, from interface{}, to interface{}) *Query_ {
	var m M
	if v, ok := o.m["range"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m["range"] = m
	}
	m[col] = M{"from": from, "to": to}
	return o
}

//Filter, func always ending with _
func (o *Query_) Equal_(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.mf["term"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf["term"] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Prefix_(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.mf["prefix"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf["prefix"] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Range_(col string, from interface{}, to interface{}) *Query_ {
	var m M
	if v, ok := o.mf["range"]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf["range"] = m
	}
	m[col] = M{"from": from, "to": to}
	return o
}

/* Result Object */
type Result_ struct {
	Data	[]*Data_
	Found	int	//total data found
	Count	int	//total data returned
}

/* Query Iterator */
type Iter_ struct {
	query		*Query_
	res		*Result_
	cur, Found	int
}

func NewIter(q *Query_) (*Iter_, error) {
	res, err := q.Find()
	if err != nil {
		return nil, err
	}
	it := &Iter_{query: q, res: res, Found: res.Found}
	return it, nil
}

func (o *Iter_) Next() *Data_ {
	if (o.cur + o.query.offset) >= o.Found {
		return nil
	}
	if o.cur >= o.res.Count {
		res, err := o.query.Offset(o.query.offset + o.cur).Find()
		if err != nil {
			return nil
		}
		if res.Found < o.Found {	//some data has been deleted
			if (o.cur + o.query.offset) >= res.Found {
				return nil
			}
			o.Found = res.Found
		}
		o.res, o.cur = res, 0
	}
	d := o.res.Data[o.cur]
	o.cur++
	return d
}

/* Helper */
func GetBodyJson(resp *http.Response) (M, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := M{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func SetBodyJson(req *http.Request, m M) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req.Body = nopCloser{bytes.NewBuffer(b)}
	return nil
}

func SetBodyString(req *http.Request, s string) {
	req.Body = nopCloser{bytes.NewBufferString(s)}
}
