package main

func getScore(w http.ResponseWriter, r *http.Request) {
  // Set proper content-type header for jsonp
  w.Header().Set("Content-Type", "text/javascript")

  callback := r.FormValue("callback")
  
  s1 := Score{"Mika", 64}
  s2 := Score{"Mikko", 62}
  s3 := Score{"Pekko", 34}
  s4 := Score{"Arimas", 95}

  var resp = struct {
     Result []Score
  }{
     Result: []Score{s1, s2, s3, s4},
  }
  b, err := json.Marshal(&resp)
  if err != nil {}
  res := callback + "(" + string(b) + ")"
  fmt.Fprint(w, res)
}
