package main
import (
	"net/http"
	"wrappers/youtube"
	"wrappers/itunes"
	"helpers/types"
	"fmt"
	"sync"
	"encoding/json"
	"runtime"
)



func queryAll(w http.ResponseWriter, r *http.Request) {
	//movieQuery := "inception"
	movieQuery := r.FormValue("keyword")

	//allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var wg sync.WaitGroup
	//services := []string{"youtube", "itunes"}

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
	results = append(results, <-ch)
	wg.Wait()
	close(ch)

	resultBlob, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error marshaling slice of MovieResult structures")
	}
	fmt.Fprintf(w, string(resultBlob))
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/query", queryAll)

	http.ListenAndServe("localhost:9999", nil)
}

