package main

func getScore(w http.ResponseWriter, r *http.Request) {
  // Set proper content-type header for jsonp
  w.Header().Set(&#34;Content-Type&#34;, &#34;text/javascript&#34;)

  callback := r.FormValue(&#34;callback&#34;)
  
  s1 := Score{&#34;Mika&#34;, 64}
  s2 := Score{&#34;Mikko&#34;, 62}
  s3 := Score{&#34;Pekko&#34;, 34}
  s4 := Score{&#34;Arimas&#34;, 95}

  var resp = struct {
     Result []Score
  }{
     Result: []Score{s1, s2, s3, s4},
  }
  b, err := json.Marshal(&amp;resp)
  if err != nil {}
  res := callback &#43; &#34;(&#34; &#43; string(b) &#43; &#34;)&#34;
  fmt.Fprint(w, res)
}
