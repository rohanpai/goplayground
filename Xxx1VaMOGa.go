package main

import (
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;os/exec&#34;
	&#34;strings&#34;

	&#34;github.com/bmizerany/mc&#34;
	&#34;github.com/elazarl/goproxy&#34;
	&#34;github.com/jackc/pgx&#34;
	&#34;github.com/pquerna/ffjson/ffjson&#34;
)

var pool *pgx.ConnPool

func UrlHasPrefix(prefix string) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.Method == http.MethodGet &amp;&amp; strings.HasPrefix(req.URL.Path, prefix) &amp;&amp; !strings.HasPrefix(req.URL.Path, &#34;/packages/search/&#34;)
	}
}

func PathIs(path string) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.Method == http.MethodGet &amp;&amp; req.URL.Path == path
	}
}

func getEnv(key, def string) string {
	k := os.Getenv(key)
	if k == &#34;&#34; {
		return def
	}
	return k
}

func main() {
	memcachedUrl := getEnv(&#34;MEMCACHEDCLOUD_SERVERS&#34;, &#34;localhost:11211&#34;)
	cn, err := mc.Dial(&#34;tcp&#34;, memcachedUrl)
	if err != nil {
		log.Fatalf(&#34;Memcached connection error: %s&#34;, err)
	}

	memcachedUsername := os.Getenv(&#34;MEMCACHEDCLOUD_USERNAME&#34;)
	memcachedPassword := os.Getenv(&#34;MEMCACHEDCLOUD_PASSWORD&#34;)
	if memcachedUrl != &#34;&#34; &amp;&amp; memcachedPassword != &#34;&#34; {
		if err := cn.Auth(memcachedUsername, memcachedPassword); err != nil {
			log.Fatalf(&#34;Memcached auth error: %s&#34;, err)
		}
	}

	pgxcfg, err := pgx.ParseURI(os.Getenv(&#34;DATABASE_URL&#34;))
	if err != nil {
		log.Fatalf(&#34;Parse URI error: %s&#34;, err)
	}
	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxcfg,
		MaxConnections: 20,
		AfterConnect: func(conn *pgx.Conn) error {
			_, err := conn.Prepare(&#34;getPackage&#34;, `SELECT name, url FROM packages WHERE name = $1`)
			return err
		},
	})
	if err != nil {
		log.Fatalf(&#34;Connection error: %s&#34;, err)
	}
	defer pool.Close()

	binary, err := exec.LookPath(&#34;node&#34;)
	if err != nil {
		log.Fatalf(&#34;Could not lookup node path: %s&#34;, err)
	}

	cmd := exec.Command(binary, &#34;--expose_gc&#34;, &#34;index.js&#34;)
	env := os.Environ()
	env = append([]string{&#34;PORT=3001&#34;}, env...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf(&#34;Could not start node: %s&#34;, err)
	}
	// TODO: Does defer even work here?
	defer func() {
		if err := cmd.Wait(); err != nil {
			log.Fatalf(&#34;Node process failed: %s&#34;, err)
		}
	}()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.NonproxyHandler = http.HandlerFunc(NonProxy)
	proxy.OnRequest(PathIs(&#34;/packages&#34;)).DoFunc(ListPackages)
	proxy.OnRequest(UrlHasPrefix(&#34;/packages/&#34;)).DoFunc(GetPackage)

	port := getEnv(&#34;PORT&#34;, &#34;3000&#34;)
	log.Println(&#34;Starting web server at port&#34;, port)
	log.Fatal(http.ListenAndServe(&#34;:&#34;&#43;port, proxy))
}

func NonProxy(w http.ResponseWriter, req *http.Request) {
	req.Host = &#34;bower.herokuapp.com&#34;
	req.URL.Scheme = &#34;http&#34;
	req.URL.Host = &#34;localhost:3001&#34;
	proxy.ServeHTTP(w, req)
}

func GetPackages(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	elements := strings.Split(r.URL.Path, &#34;/&#34;)
	packageName := elements[len(elements)-1]

	var name, url string
	if err := pool.QueryRow(&#34;getPackage&#34;, packageName).Scan(&amp;name, &amp;url); err != nil {
		if err == pgx.ErrNoRows {
			return r, goproxy.NewResponse(r, &#34;text/html&#34;, http.StatusNotFound, &#34;Package not found&#34;)
		}
		return r, goproxy.NewResponse(r, &#34;text/html&#34;, http.StatusInternalServerError, &#34;Internal server error&#34;)
	}

	// TODO: Why use a map[string]string instead of a struct?
	result := map[string]string{&#34;name&#34;: name, &#34;url&#34;: url}
	resultByteArray, _ := ffjson.Marshal(result)
	return r, goproxy.NewResponse(r, &#34;application/json&#34;, http.StatusOK, string(resultByteArray))
}

func ListPackages(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	val, _, _, err := cn.Get(&#34;packages&#34;)
	if err != nil {
		return r, nil
	}
	log.Println(&#34;MEMCACHED FROM GO SERVER&#34;)
	return r, goproxy.NewResponse(r, &#34;application/json&#34;, http.StatusOK, val)
}
