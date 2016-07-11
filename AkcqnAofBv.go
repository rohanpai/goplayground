package main

import (
	&#34;fmt&#34;
	&#34;github.com/sourcegraph/go-webkit2/webkit2&#34;
	&#34;github.com/sqs/gojs&#34;
	&#34;github.com/sqs/gotk3/glib&#34;
	&#34;github.com/sqs/gotk3/gtk&#34;
	&#34;runtime&#34;
)

func main() {
	runtime.LockOSThread()
	gtk.Init(nil)

	webView := webkit2.NewWebView()
	defer webView.Destroy()

	webView.Connect(&#34;load-failed&#34;, func() {
		fmt.Println(&#34;Load failed.&#34;)
	})
	webView.Connect(&#34;load-changed&#34;, func(ctx *glib.CallbackContext) {
		loadEvent := webkit2.LoadEvent(ctx.Arg(0).Int())
		switch loadEvent {
		case webkit2.LoadFinished:
			fmt.Println(&#34;Load finished.&#34;)
			fmt.Printf(&#34;Title: %q\n&#34;, webView.Title())
			fmt.Printf(&#34;URI: %s\n&#34;, webView.URI())
			webView.RunJavaScript(&#34;window.location.hostname&#34;, func(val *gojs.Value, err error) {
				if err != nil {
					fmt.Println(&#34;JavaScript error.&#34;)
				} else {
					fmt.Printf(&#34;Hostname (from JavaScript): %q\n&#34;, val)
				}
				gtk.MainQuit()
			})
		}
	})

	glib.IdleAdd(func() bool {
		webView.LoadURI(&#34;https://www.google.com/&#34;)
		return false
	})

	gtk.Main()

}
