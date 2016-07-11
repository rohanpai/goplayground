package main

import (
	&#34;fmt&#34;
	&#34;github.com/Kaey/framebuffer&#34;
	&#34;log&#34;
)

func main() {
	fb, err := framebuffer.Init(&#34;/dev/fb0&#34;)
	if err != nil {
		log.Fatalln(err)
	}
	defer fb.Close()
	fb.Clear(0, 0, 0, 0)
	fb.WritePixel(200, 100, 255, 0, 0, 0)
	fmt.Scanln()
}
