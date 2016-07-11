package main

import (
	&#34;crypto/md5&#34;
	&#34;encoding/hex&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net/http&#34;
	&#34;net/url&#34;
	&#34;strings&#34;
)

type Authorization struct {
	Username, Password, Realm, NONCE, QOP, Opaque, Algorithm string
}

func GetAuthorization(username, password string, resp *http.Response) *Authorization {
	header := resp.Header.Get(&#34;www-authenticate&#34;)
	parts := strings.SplitN(header, &#34; &#34;, 2)
	parts = strings.Split(parts[1], &#34;, &#34;)
	fmt.Println(&#34;Parts: &#34;, parts)
	opts := make(map[string]string)

	for _, part := range parts {
		fmt.Println(&#34;Part: &#34;, part)
		vals := strings.SplitN(part, &#34;=&#34;, 2)
		key := vals[0]
		val := strings.Trim(vals[1], &#34;\&#34;,&#34;)
		opts[key] = val
	}

	auth := Authorization{
		username, password,
		opts[&#34;realm&#34;], opts[&#34;nonce&#34;], opts[&#34;qop&#34;], opts[&#34;opaque&#34;], opts[&#34;algorithm&#34;],
	}
	return &amp;auth
}

func SetDigestAuth(r *http.Request, username, password string, resp *http.Response, nc int) {
	auth := GetAuthorization(username, password, resp)
	auth_str := GetAuthString(auth, r.URL, r.Method, nc)
	r.Header.Add(&#34;Authorization&#34;, auth_str)
}

func GetAuthString(auth *Authorization, url *url.URL, method string, nc int) string {
	a1 := auth.Username &#43; &#34;:&#34; &#43; auth.Realm &#43; &#34;:&#34; &#43; auth.Password
	h := md5.New()
	io.WriteString(h, a1)
	ha1 := hex.EncodeToString(h.Sum(nil))

	h = md5.New()
	a2 := method &#43; &#34;:&#34; &#43; url.Path
	io.WriteString(h, a2)
	ha2 := hex.EncodeToString(h.Sum(nil))

	nc_str := fmt.Sprintf(&#34;%08x&#34;, nc)
	hnc := &#34;MTM3MDgw&#34;
	
	respdig := fmt.Sprintf(&#34;%s:%s:%s:%s:%s:%s&#34;, ha1, auth.NONCE, nc_str, hnc, auth.QOP, ha2)
	h = md5.New()
	io.WriteString(h, respdig)
	respdig = hex.EncodeToString(h.Sum(nil))

	base := &#34;username=\&#34;%s\&#34;, realm=\&#34;%s\&#34;, nonce=\&#34;%s\&#34;, uri=\&#34;%s\&#34;, response=\&#34;%s\&#34;&#34;
	base = fmt.Sprintf(base, auth.Username, auth.Realm, auth.NONCE, url.Path, respdig)
	if auth.Opaque != &#34;&#34; {
		base &#43;= fmt.Sprintf(&#34;, opaque=\&#34;%s\&#34;&#34;, auth.Opaque)
	}
	if auth.QOP != &#34;&#34; {
		base &#43;= fmt.Sprintf(&#34;, qop=\&#34;%s\&#34;, nc=%s, cnonce=\&#34;%s\&#34;&#34;, auth.QOP, nc_str, hnc)
	}
	if auth.Algorithm != &#34;&#34; {
		base &#43;= fmt.Sprintf(&#34;, algorithm=\&#34;%s\&#34;&#34;, auth.Algorithm)
	}

	// r.Header.Add(&#34;Authorization&#34;, &#34;Digest &#34; &#43;base)
	return &#34;Digest &#34; &#43; base
}

func main() {
	auth := Authorization{
		&#34;Joe&#34;, &#34;Schmoe&#34;, &#34;RETS Server&#34;, 
		&#34;0decbedc0f7828a0f8a0f0ea4a2107e3&#34;, &#34;auth&#34;, 
		&#34;a0531450f6a92cda6b30ae8e4de1cd2b&#34;, &#34;&#34;,
	}
	myURL := url.URL{&#34;http&#34;, &#34;&#34;, nil, &#34;www.dis.com:6103&#34;, &#34;/rets/login&#34;, &#34;&#34;, &#34;&#34;}
	fmt.Println(GetAuthString(&amp;auth, &amp;myURL, &#34;POST&#34;, 3))
}