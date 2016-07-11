package main

import (
	"github.com/robertkrimen/otto"
	"log"
)

func main() {
	log.Printf("Creating JS interpreter")
	js := otto.New()

	var function otto.Value

	log.Printf("Defining setEnrichFunction")
	js.Set("setEnrichFunction", func(call otto.FunctionCall) otto.Value {
		function = call.Argument(0)
		if class := function.Class(); class != "Function" {
			log.Fatalf("setEnrichFunction: expected Function, got %s instead.", class)
		}
		return otto.UndefinedValue()
	})

	log.Printf("Registering enrich function")
	js.Run(`
		setEnrichFunction(function(data) {
			data.timestamp = new Date().toUTCString();
		});
	`)

	data := map[string]string{
		"foo": "bar",
		"theAnswer": "42",
	}

	log.Printf("raw data: %#v", data)

	arg, err := js.ToValue(data)
	if err != nil {
		log.Fatalf("couldn't convert message to JS value")
	}

	log.Printf("Calling enrich function")
	_, err = function.Call(otto.NullValue(), arg)
	if err != nil {
		log.Fatalf("calling enrich function failed: %v", err)
	}

	log.Printf("enriched data: %#v", data)
}