package main

import (
	&#34;encoding/json&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;net/http&#34;
	&#34;os&#34;
)

//Coordinates is a set of coordinate
type Coordinates []Coordinate

//Coordinate is a [longitude, latitude]
type Coordinate [2]float64

// Point rapresent a geojson point geometry object
type Point struct {
	Type       string `json:&#34;type&#34;`
	Coordinate `json:&#34;coordinates&#34;`
}

func (p *Point) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&amp;p)
	return err
}

// Polygon rapresent a geojson polygon geometry object
type Polygon struct {
	Type        string `json:&#34;type&#34;`
	Coordinates `json:&#34;coordinates&#34;`
}

func (p *Polygon) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&amp;p)
	return err
}

// MultiPolygon rapresent a geojson mulitpolygon  geometry object
type MultiPolygon struct {
	Type        string        `json:&#34;type&#34;`
	Coordinates []Coordinates `json:&#34;coordinates&#34;`
}

func (p *MultiPolygon) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&amp;p)
	return err
}

type geojson interface {
	parseJSON(io.Reader) error
}

// Endpoint is  the name of the geojson handler endpoint
const Endpoint = &#34;/tos2/geojson/&#34;

//Handler handles a request for a geojsonPoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// request
	resp, err := matcher(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// response
	encoder := json.NewEncoder(w)
	w.Header().Set(&#34;Content-Type&#34;, &#34;application/json&#34;)
	err = encoder.Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Matcher exctract from the url witch geoJSON object we want
func matcher(r *http.Request) (resp geojson, err error) {
	objectType := r.URL.Path[len(Endpoint):]
	switch objectType {
	case &#34;polygon&#34;:
		p := Polygon{}
		err = p.parseJSON(r.Body)
		return &amp;p, err
	case &#34;point&#34;:
		p := Point{}
		err = p.parseJSON(r.Body)
		return &amp;p, err
	case &#34;multipolygon&#34;:
		p := MultiPolygon{}
		err = p.parseJSON(r.Body)
		return &amp;p, err
	default:
		err = fmt.Errorf(&#34;Bad geoJSON object type&#34;)
	}
	return
}

func main() {
	// geojson
	geoJSONHandler := http.HandlerFunc(Handler)
	http.Handle(Endpoint, geoJSONHandler)
	// server
	port := os.Getenv(&#34;PORT&#34;)
	addr := fmt.Sprintf(&#34;:%v&#34;, port)
	http.ListenAndServe(addr, nil)
}