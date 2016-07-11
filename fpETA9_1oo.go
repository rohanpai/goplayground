package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

var oldButtonFormat = regexp.MustCompile("^/([0-9]*)/([0-9]).(jpg|gif|png)$")

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	slices := oldButtonFormat.FindStringSubmatch(r.URL.Path)
	if len(slices) == 4 {
		// Display old images
		w.Header().Set("Cache-Control", "max-age:290304000, public")
		w.Header().Set("Last-Modified", cacheSince)
		w.Header().Set("Expires", cacheUntil)

		switch slices[3] {
		case "gif":
			w.Header().Set("Content-Type", "image/gif")
			break
		case "png":
			w.Header().Set("Content-Type", "image/png")
			break
		case "jpg":
			w.Header().Set("Content-Type", "image/jpeg")
			break
		}

		w.Write(oldButtons[slices[3]])
	} else {
		// Display standard paths
		switch r.URL.Path {
		case "/":
			http.Redirect(w, r, "http://comicrank.com", 302)
			break
		case "/robots.txt":
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "User-agent: *\nDisallow: /")
			break
		default:
			http.NotFound(w, r)
			break
		}
	}
}

// Map to hold old button images in memory
var oldButtons map[string][]byte

// Cache given button image's data in memory
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
	initImg("gif", "gif.gif")
	initImg("jpg", "jpg.jpg")
	initImg("png", "png.png")

	// Start the http server
	http.HandleFunc("/", baseHandler)
	http.ListenAndServe(":8080", nil)
}
