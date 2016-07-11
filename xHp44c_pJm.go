package main
import (
	&#34;net/http&#34;
	&#34;wrappers/youtube&#34;
	&#34;wrappers/itunes&#34;
	&#34;helpers/types&#34;
	&#34;fmt&#34;
	&#34;sync&#34;
	&#34;encoding/json&#34;
	&#34;runtime&#34;
)



func queryAll(w http.ResponseWriter, r *http.Request) {
	//movieQuery := &#34;inception&#34;
	movieQuery := r.FormValue(&#34;keyword&#34;)

	//allow cross domain AJAX requests
	w.Header().Set(&#34;Access-Control-Allow-Origin&#34;, &#34;*&#34;)

	var wg sync.WaitGroup
	//services := []string{&#34;youtube&#34;, &#34;itunes&#34;}

	ch := make(chan types.MovieResult, 2)
	results := []types.MovieResult{}

	wg.Add(1)
	go func() {
		youtube.Query(movieQuery, ch)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		itunes.Query(movieQuery, ch)
		wg.Done()
	}()
	results = append(results, &lt;-ch)
	wg.Wait()
	close(ch)

	resultBlob, err := json.Marshal(results)
	if err != nil {
		fmt.Println(&#34;Error marshaling slice of MovieResult structures&#34;)
	}
	fmt.Fprintf(w, string(resultBlob))
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc(&#34;/query&#34;, queryAll)

	http.ListenAndServe(&#34;localhost:9999&#34;, nil)
}

