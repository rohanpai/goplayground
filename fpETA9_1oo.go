package main

import (
	&#34;fmt&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;regexp&#34;
	&#34;time&#34;
)

var oldButtonFormat = regexp.MustCompile(&#34;^/([0-9]*)/([0-9]).(jpg|gif|png)$&#34;)

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	slices := oldButtonFormat.FindStringSubmatch(r.URL.Path)
	if len(slices) == 4 {
		// Display old images
		w.Header().Set(&#34;Cache-Control&#34;, &#34;max-age:290304000, public&#34;)
		w.Header().Set(&#34;Last-Modified&#34;, cacheSince)
		w.Header().Set(&#34;Expires&#34;, cacheUntil)

		switch slices[3] {
		case &#34;gif&#34;:
			w.Header().Set(&#34;Content-Type&#34;, &#34;image/gif&#34;)
			break
		case &#34;png&#34;:
			w.Header().Set(&#34;Content-Type&#34;, &#34;image/png&#34;)
			break
		case &#34;jpg&#34;:
			w.Header().Set(&#34;Content-Type&#34;, &#34;image/jpeg&#34;)
			break
		}

		w.Write(oldButtons[slices[3]])
	} else {
		// Display standard paths
		switch r.URL.Path {
		case &#34;/&#34;:
			http.Redirect(w, r, &#34;http://comicrank.com&#34;, 302)
			break
		case &#34;/robots.txt&#34;:
			w.Header().Set(&#34;Content-Type&#34;, &#34;text/plain&#34;)
			fmt.Fprintf(w, &#34;User-agent: *\nDisallow: /&#34;)
			break
		default:
			http.NotFound(w, r)
			break
		}
	}
}

// Map to hold old button images in memory
var oldButtons map[string][]byte

// Cache given button image&#39;s data in memory
func initImg(index string, filename string) {
	file, _ := os.Open(filename)
	info, _ := file.Stat()
	oldButtons[index] = make([]byte, info.Size())
	file.Read(oldButtons[index])
	file.Close()
}

func main() {
	// Cache old buttons in memory
	oldButtons = make(map[string][]byte, 3)
	initImg(&#34;gif&#34;, &#34;gif.gif&#34;)
	initImg(&#34;jpg&#34;, &#34;jpg.jpg&#34;)
	initImg(&#34;png&#34;, &#34;png.png&#34;)

	// Start the http server
	http.HandleFunc(&#34;/&#34;, baseHandler)
	http.ListenAndServe(&#34;:8080&#34;, nil)
}
