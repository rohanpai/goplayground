package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/crypto/scrypt"
)

type Status struct {
	Code   int               `json:"code"`
	Name   string            `json:"name"`
	Desc   string            `json:"desc"`
	Fields map[string]string `json:"fields"`
}

type GetSaltResult struct {
	Status       `json:"status"`
	Salt         string `json:"salt"`
	CSRFToken    string `json:"csrf_token"`
	LoginSession string `json:"login_session"`
}

func main() {
	var urlString string
	var url *url.URL
	email := "MyEmail"
	password := []byte("MyPassword")

	// Get Salt
	urlString = "https://keybase.io/_/api/1.0/getsalt.json"
	url, _ = url.Parse(urlString)
	q := url.Query()
	q.Set("email_or_username", email)
	url.RawQuery = q.Encode()
	resp, _ := http.Get(url.String())
	result := &GetSaltResult{}
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
	urlString = "https://keybase.io/_/api/1.0/login.json"
	url, _ = url.Parse(urlString)
	q = url.Query()
	q.Set("email_or_username", email)
	q.Set("hmac_pwh", hex.EncodeToString(mac.Sum(nil)))
	q.Set("login_session", base64.StdEncoding.EncodeToString(loginSession))
	q.Set("csrf_token", result.CSRFToken)
	resp, _ = http.PostForm(url.String(), q)
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(contents))
}