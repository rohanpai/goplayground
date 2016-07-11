package main

import (
	&#34;github.com/robertkrimen/otto&#34;
	&#34;log&#34;
)

func main() {
	log.Printf(&#34;Creating JS interpreter&#34;)
	js := otto.New()

	var function otto.Value

	log.Printf(&#34;Defining setEnrichFunction&#34;)
	js.Set(&#34;setEnrichFunction&#34;, func(call otto.FunctionCall) otto.Value {
		function = call.Argument(0)
		if class := function.Class(); class != &#34;Function&#34; {
			log.Fatalf(&#34;setEnrichFunction: expected Function, got %s instead.&#34;, class)
		}
		return otto.UndefinedValue()
	})

	log.Printf(&#34;Registering enrich function&#34;)
	js.Run(`
		setEnrichFunction(function(data) {
			data.timestamp = new Date().toUTCString();
		});
	`)

	data := map[string]string{
		&#34;foo&#34;: &#34;bar&#34;,
		&#34;theAnswer&#34;: &#34;42&#34;,
	}

	log.Printf(&#34;raw data: %#v&#34;, data)

	arg, err := js.ToValue(data)
	if err != nil {
		log.Fatalf(&#34;couldn&#39;t convert message to JS value&#34;)
	}

	log.Printf(&#34;Calling enrich function&#34;)
	_, err = function.Call(otto.NullValue(), arg)
	if err != nil {
		log.Fatalf(&#34;calling enrich function failed: %v&#34;, err)
	}

	log.Printf(&#34;enriched data: %#v&#34;, data)
}