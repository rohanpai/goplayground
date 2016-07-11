package main

import (
	&#34;crypto/hmac&#34;
	&#34;crypto/sha512&#34;
	&#34;encoding/base64&#34;
	&#34;encoding/hex&#34;
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;io/ioutil&#34;
	&#34;net/http&#34;
	&#34;net/url&#34;

	&#34;golang.org/x/crypto/scrypt&#34;
)

type Status struct {
	Code   int               `json:&#34;code&#34;`
	Name   string            `json:&#34;name&#34;`
	Desc   string            `json:&#34;desc&#34;`
	Fields map[string]string `json:&#34;fields&#34;`
}

type GetSaltResult struct {
	Status       `json:&#34;status&#34;`
	Salt         string `json:&#34;salt&#34;`
	CSRFToken    string `json:&#34;csrf_token&#34;`
	LoginSession string `json:&#34;login_session&#34;`
}

func main() {
	var urlString string
	var url *url.URL
	email := &#34;MyEmail&#34;
	password := []byte(&#34;MyPassword&#34;)

	// Get Salt
	urlString = &#34;https://keybase.io/_/api/1.0/getsalt.json&#34;
	url, _ = url.Parse(urlString)
	q := url.Query()
	q.Set(&#34;email_or_username&#34;, email)
	url.RawQuery = q.Encode()
	resp, _ := http.Get(url.String())
	result := &amp;GetSaltResult{}
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(result)

	// crypto
	salt, _ := hex.DecodeString(result.Salt)
	loginSession, _ := base64.StdEncoding.DecodeString(result.LoginSession)
	skey, _ := scrypt.Key(password, salt, 32768, 8, 1, 224)
	pwh := skey[192:224]
	mac := hmac.New(sha512.New, loginSession)
	mac.Write(pwh)

	// Login
	urlString = &#34;https://keybase.io/_/api/1.0/login.json&#34;
	url, _ = url.Parse(urlString)
	q = url.Query()
	q.Set(&#34;email_or_username&#34;, email)
	q.Set(&#34;hmac_pwh&#34;, hex.EncodeToString(mac.Sum(nil)))
	q.Set(&#34;login_session&#34;, base64.StdEncoding.EncodeToString(loginSession))
	q.Set(&#34;csrf_token&#34;, result.CSRFToken)
	resp, _ = http.PostForm(url.String(), q)
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(contents))
}