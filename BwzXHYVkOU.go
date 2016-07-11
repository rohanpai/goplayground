package db

import (
	&#34;bytes&#34;
	&#34;encoding/json&#34;
	&#34;errors&#34;
	&#34;flag&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;net/url&#34;
	&#34;reflect&#34;
	&#34;time&#34;
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
	TypeGeo		= 9	//uses string format &#34;lat,lon&#34;

	DefQueryLimit		= 10
	EsHLFragment		= 1
	EsHLFragmentSize	= 100

	DateFormat	= &#34;2006-01-02&#34;
	EsDateFormat	= &#34;YYYY-MM-dd&#34;
)

type M map[string]interface{}
type nopCloser struct{ io.Reader }

func (nopCloser) Close() error	{ return nil }

var (
	Remap	bool	//re-map Riak buckets to reflect definition, (ES not supported, must create &#43; reindex)

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
	flag.BoolVar(&amp;Remap, &#34;dbremap&#34;, false, &#34;Remap DB to reflect changes in definition&#34;)

	Host = &#34;127.0.0.1&#34;
	Port = &#34;8098&#34;
	EsHost = &#34;127.0.0.1&#34;
	EsPort = &#34;9200&#34;

	EsHLTagPre, EsHLTagPost = []string{&#34;&lt;strong&gt;&#34;}, []string{&#34;&lt;/strong&gt;&#34;}
}

func Init(cfgFile string) error {
	//load config
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	m := M{}
	err = json.Unmarshal(b, &amp;m)
	if err != nil {
		return err
	}

	if v, ok := m[&#34;name&#34;]; ok {
		Name = v.(string)
	}
	if v, ok := m[&#34;host&#34;]; ok {
		Host = v.(string)
	}
	if v, ok := m[&#34;port&#34;]; ok {
		Port = v.(string)
	}
	if v, ok := m[&#34;es_host&#34;]; ok {
		EsHost = v.(string)
	}
	if v, ok := m[&#34;es_port&#34;]; ok {
		EsPort = v.(string)
	}

	tables = make(map[string]*table_)
	if v, ok := m[&#34;tables&#34;]; ok {
		tablesM := v.([]interface{})
		for _, v := range tablesM {
			t := v.(map[string]interface{})
			table := &amp;table_{Name: t[&#34;name&#34;].(string), Cols: make(map[string]*tableCol_)}
			if v, ok := t[&#34;key&#34;]; ok {
				table.Key = v.(string)
			}
			if v, ok := t[&#34;es&#34;]; ok {
				table.Es = v.(bool)
			}
			if v, ok := t[&#34;cols&#34;]; ok {	//columns
				tc := v.([]interface{})
				for _, v := range tc {
					t := v.(map[string]interface{})
					col := &amp;tableCol_{Name: t[&#34;name&#34;].(string)}
					switch t[&#34;type&#34;].(string) {
					case &#34;string&#34;:
						col.Type = TypeString
					case &#34;bool&#34;:
						col.Type = TypeBool
					case &#34;int&#34;:
						col.Type = TypeInt
					case &#34;int32&#34;:
						col.Type = TypeInt32
					case &#34;int16&#34;:
						col.Type = TypeInt16
					case &#34;float&#34;:
						col.Type = TypeFloat
					case &#34;date&#34;:
						col.Type = TypeDate
					case &#34;time&#34;:
						col.Type = TypeTime
					case &#34;geo&#34;:
						col.Type = TypeGeo
					}

					if v, ok := t[&#34;index&#34;]; ok {
						col.Index = v.(string)
					}
					if v, ok := t[&#34;boost&#34;]; ok {
						col.Boost = v.(float64)
					}
					table.Cols[col.Name] = col
				}
			}
			tables[table.Name] = table
		}
	}

	//connections
	ConnStr = &#34;http://&#34; &#43; Host &#43; &#34;:&#34; &#43; Port &#43; &#34;/&#34;
	EsConnStr = &#34;http://&#34; &#43; EsHost &#43; &#34;:&#34; &#43; EsPort &#43; &#34;/&#34;

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
	//riak -&gt; GET /ping
	req, err := http.NewRequest(&#34;GET&#34;, ConnStr&#43;&#34;ping&#34;, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(&#34;Riak not ready&#34;)
	}

	//ES -&gt; GET /_cluster/health
	req, err = http.NewRequest(&#34;GET&#34;, EsConnStr, nil)
	if err != nil {
		return err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(&#34;ES request failed&#34;)
	}
	m, err := GetBodyJson(resp)
	if err != nil {
		return err
	}
	if !m[&#34;ok&#34;].(bool) {
		return errors.New(&#34;ES not ready&#34;)
	}
	return nil
}

/* Table */
func OpenTable(table string) error {
	tbl, ok := tables[table]
	if !ok {
		return errors.New(&#34;Table not defined&#34;)
	}
	log.Println(&#34;Checking table &#34; &#43; table)

	//riak -&gt; GET /buckets/bucket/props
	req, err := http.NewRequest(&#34;GET&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;table&#43;&#34;/props&#34;, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if Remap || resp.StatusCode != 200 {	//create it
		//riak -&gt; PUT /buckets/bucket/props
		req, err = http.NewRequest(&#34;PUT&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;table&#43;&#34;/props&#34;, nil)
		if err != nil {
			return err
		}
		req.Header.Set(&#34;Content-Type&#34;, &#34;application/json&#34;)
		SetBodyString(req, `{&#34;props&#34;:{&#34;allow_mult&#34;:false, &#34;last_write_wins&#34;:true}}`)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 204 {
			return errors.New(&#34;Riak fail to setup bucket&#34;)
		}
		log.Println(&#34;Riak bucket props for &#34; &#43; table &#43; &#34; created&#34;)
	}

	if tbl.Es {
		//ES -&gt; GET /index/_mapping
		req, err = http.NewRequest(&#34;GET&#34;, EsConnStr&#43;table&#43;&#34;/_mapping&#34;, nil)
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
			log.Println(&#34;ES mapping for &#34; &#43; table &#43; &#34; created&#34;)
		}
	}

	log.Println(&#34;Table &#34; &#43; table &#43; &#34; OK&#34;)
	return nil
}

func esMapTable(table *table_) error {
	mc := M{}
	for _, v := range table.Cols {
		c := M{&#34;store&#34;: &#34;yes&#34;}	//store everything since we&#39;re using this as main DB
		switch v.Type {
		case TypeString:
			c[&#34;type&#34;] = &#34;string&#34;
		case TypeBool:
			c[&#34;type&#34;] = &#34;boolean&#34;
		case TypeFloat:
			c[&#34;type&#34;] = &#34;double&#34;
		case TypeInt:
			c[&#34;type&#34;] = &#34;long&#34;
		case TypeInt32:
			c[&#34;type&#34;] = &#34;integer&#34;
		case TypeInt16:
			c[&#34;type&#34;] = &#34;short&#34;
		case TypeDate:
			c[&#34;type&#34;] = &#34;date&#34;
			c[&#34;format&#34;] = EsDateFormat
		case TypeTime:
			c[&#34;type&#34;] = &#34;date&#34;	//use ES default dateoptionaltime RFC3339
		case TypeGeo:
			c[&#34;type&#34;] = &#34;geo_point&#34;
			c[&#34;lat_lon&#34;] = true
		}
		if v.Index == &#34;&#34; || v.Index == &#34;no&#34; {
			c[&#34;index&#34;] = &#34;no&#34;
		} else if v.Index == &#34;yes&#34; {
			//ES default to indexed
		} else if v.Index == &#34;exact&#34; {
			c[&#34;index&#34;] = &#34;not_analyzed&#34;
		} else {	//fulltext, etc.
			c[&#34;index&#34;] = &#34;analyzed&#34;
			c[&#34;analyzer&#34;] = &#34;ih_&#34; &#43; v.Index
		}

		if v.Boost &gt; 0 {
			c[&#34;boost&#34;] = v.Boost
		}
		mc[v.Name] = c
	}
	m := M{&#34;doc&#34;: M{&#34;ignore_conflicts&#34;: false, &#34;properties&#34;: mc}}
	m = M{&#34;mappings&#34;: m}

	req, err := http.NewRequest(&#34;PUT&#34;, EsConnStr&#43;table.Name, nil)
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
		return errors.New(&#34;ES mapping fail&#34;)
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
	return &amp;Data_{table: tables[table], m: make(M)}
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
	return &#34;&#34;
}

func (o *Data_) Set(col string, v interface{}) *Data_ {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.String {	//add only if not empty
		if r.Len() &gt; 0 {
			o.m[col] = v
		}
	}
	return o
}
func (o *Data_) SetStr(col string, v string) *Data_ {
	if len(v) &gt; 0 {
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
	if len(v) &gt; 0 {
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
	return &#34;&#34;
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
	return &#34;&#34;
}
func (o *Data_) Highlight(col string) string {
	if o.mh != nil {
		v, ok := o.mh[col]
		if ok {
			v := v.([]interface{})
			return v[0].(string)
		}
	}
	return &#34;&#34;
}

//Operations
func (o *Data_) Put() error	{ return o.put(false) }
func (o *Data_) CondPut() error	{ return o.put(true) }	//save only if not exists

func (o *Data_) put(cond bool) error {
	v, ok := o.m[o.table.Key]
	if !ok {
		return errors.New(&#34;Key is not defined&#34;)
	}

	if cond {	//only put if none
		err := o.Exists(v.(string))
		if err != nil {
			return err
		}
		if o.Found {
			return errors.New(&#34;Key already existed&#34;)
		}
	}

	id := url.QueryEscape(v.(string))
	j, err := json.Marshal(o.m)
	if err != nil {
		return err
	}

	//riak -&gt; PUT /buckets/bucket/keys/key
	req, err := http.NewRequest(&#34;PUT&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;o.table.Name&#43;&#34;/keys/&#34;&#43;id, nil)
	if err != nil {
		return err
	}
	req.Header.Set(&#34;Content-Type&#34;, &#34;application/json&#34;)
	req.Body = nopCloser{bytes.NewBuffer(j)}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 &amp;&amp; resp.StatusCode != 204 {
		return errors.New(&#34;Riak store fail&#34;)
	}

	if o.table.Es {
		//ES -&gt; PUT /tablename/doc/id
		req, err = http.NewRequest(&#34;PUT&#34;, EsConnStr&#43;o.table.Name&#43;&#34;/doc/&#34;&#43;id, nil)
		if err != nil {
			return err
		}
		req.Body = nopCloser{bytes.NewBuffer(j)}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 &amp;&amp; resp.StatusCode != 201 {
			return errors.New(&#34;ES index failed&#34;)
		}
	}

	return nil
}

func (o *Data_) Exists(key string) error {	//Riak only
	//GET /buckets/bucket/keys/key
	req, err := http.NewRequest(&#34;GET&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;o.table.Name&#43;&#34;/keys/&#34;&#43;url.QueryEscape(key), nil)
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
		return errors.New(&#34;Riak fetch fail&#34;)
	}
	return nil
}

func (o *Data_) Get(key string, cols []string) error {	//Riak only, cols=nil means fetch all cols
	//GET /buckets/bucket/keys/key
	req, err := http.NewRequest(&#34;GET&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;o.table.Name&#43;&#34;/keys/&#34;&#43;url.QueryEscape(key), nil)
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
		return errors.New(&#34;Riak fetch fail&#34;)
	}
	return nil
}

func (o *Data_) Delete(key string) error {
	id := url.QueryEscape(key)

	//riak -&gt; DELETE /buckets/bucket/keys/key
	req, err := http.NewRequest(&#34;DELETE&#34;, ConnStr&#43;&#34;buckets/&#34;&#43;o.table.Name&#43;&#34;/keys/&#34;&#43;id, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 &amp;&amp; resp.StatusCode != 404 {
		return errors.New(&#34;Riak delete fail&#34;)
	}

	if o.table.Es {
		//ES -&gt; DELETE /tablename/doc/id
		req, err = http.NewRequest(&#34;DELETE&#34;, EsConnStr&#43;o.table.Name&#43;&#34;/doc/&#34;&#43;id, nil)
		if err != nil {
			return err
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 &amp;&amp; resp.StatusCode != 404 {
			return errors.New(&#34;ES delete fail&#34;)
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
	return &amp;Query_{table: tables[table], m: M{}, mf: M{}, limit: DefQueryLimit}
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
func (o *Query_) SortScore() *Query_	{ return o.Sort(&#34;_score&#34;, &#34;desc&#34;) }

func (o *Query_) Highlight(col string, size int) *Query_ {
	if o.mh == nil {
		o.mh = M{}
	}
	if size == 0 {
		o.mh[col] = M{}
	} else {
		o.mh[col] = M{&#34;fragment_size&#34;: size}
	}
	return o
}

func (o *Query_) Find() (*Result_, error) {	//ES only
	//construct query
	var m M
	if len(o.mf) &gt; 0 {
		m = M{&#34;filtered&#34;: M{&#34;query&#34;: o.m, &#34;filter&#34;: o.mf}}
	} else {
		m = o.m
	}
	m = M{&#34;query&#34;: m, &#34;from&#34;: o.offset, &#34;size&#34;: o.limit}
	if len(o.cols) &gt; 0 {
		m[&#34;fields&#34;] = o.cols
	}
	if o.ms != nil {
		m[&#34;sort&#34;] = o.ms
	}
	if o.mh != nil {
		m[&#34;highlight&#34;] = M{&#34;pre_tags&#34;:	EsHLTagPre, &#34;post_tags&#34;:	EsHLTagPost,
			&#34;number_of_fragments&#34;:	EsHLFragment, &#34;fragment_size&#34;:	EsHLFragmentSize, &#34;fields&#34;:	o.mh}
	}

	req, err := http.NewRequest(&#34;GET&#34;, EsConnStr&#43;o.table.Name&#43;&#34;/doc/_search&#34;, nil)
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
		return nil, errors.New(&#34;ES find fail&#34;)
	}
	m, err = GetBodyJson(resp)
	if err != nil {
		return nil, err
	}

	//compile results
	res := &amp;Result_{}
	m = m[&#34;hits&#34;].(map[string]interface{})
	res.Found = int(m[&#34;total&#34;].(float64))
	hits := m[&#34;hits&#34;].([]interface{})
	res.Count = len(hits)
	res.Data = make([]*Data_, res.Count)
	for i := 0; i &lt; len(hits); i&#43;&#43; {
		m = hits[i].(map[string]interface{})
		d := &amp;Data_{table: o.table, m: m[&#34;_source&#34;].(map[string]interface{}), Found: true}
		if v, ok := m[&#34;highlight&#34;]; ok {
			d.mh = v.(map[string]interface{})
		}
		res.Data[i] = d
	}

	return res, nil
}

//Collector
func (o *Query_) All() *Query_	{ o.m[&#34;match_all&#34;] = M{}; return o }

func (o *Query_) Equal(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m[&#34;term&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m[&#34;term&#34;] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Fulltext(col string, val interface{}, isAnd bool) *Query_ {
	var m M
	if v, ok := o.m[&#34;match&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m[&#34;match&#34;] = m
	}
	if isAnd {
		m[col] = M{&#34;query&#34;: val, &#34;operator&#34;: &#34;and&#34;}
	} else {
		m[col] = val
	}
	return o
}

func (o *Query_) Prefix(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m[&#34;prefix&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m[&#34;prefix&#34;] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Wildcard(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.m[&#34;wildcard&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m[&#34;wildcard&#34;] = m
	}
	m[col] = &#34;*&#34; &#43; val.(string) &#43; &#34;*&#34;
	return o
}

func (o *Query_) Range(col string, from interface{}, to interface{}) *Query_ {
	var m M
	if v, ok := o.m[&#34;range&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.m[&#34;range&#34;] = m
	}
	m[col] = M{&#34;from&#34;: from, &#34;to&#34;: to}
	return o
}

//Filter, func always ending with _
func (o *Query_) Equal_(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.mf[&#34;term&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf[&#34;term&#34;] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Prefix_(col string, val interface{}) *Query_ {
	var m M
	if v, ok := o.mf[&#34;prefix&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf[&#34;prefix&#34;] = m
	}
	m[col] = val
	return o
}

func (o *Query_) Range_(col string, from interface{}, to interface{}) *Query_ {
	var m M
	if v, ok := o.mf[&#34;range&#34;]; ok {
		m = v.(M)
	} else {
		m = M{}
		o.mf[&#34;range&#34;] = m
	}
	m[col] = M{&#34;from&#34;: from, &#34;to&#34;: to}
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
	it := &amp;Iter_{query: q, res: res, Found: res.Found}
	return it, nil
}

func (o *Iter_) Next() *Data_ {
	if (o.cur &#43; o.query.offset) &gt;= o.Found {
		return nil
	}
	if o.cur &gt;= o.res.Count {
		res, err := o.query.Offset(o.query.offset &#43; o.cur).Find()
		if err != nil {
			return nil
		}
		if res.Found &lt; o.Found {	//some data has been deleted
			if (o.cur &#43; o.query.offset) &gt;= res.Found {
				return nil
			}
			o.Found = res.Found
		}
		o.res, o.cur = res, 0
	}
	d := o.res.Data[o.cur]
	o.cur&#43;&#43;
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
	err = json.Unmarshal(body, &amp;m)
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
