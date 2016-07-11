package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

//Coordinates is a set of coordinate
type Coordinates []Coordinate

//Coordinate is a [longitude, latitude]
type Coordinate [2]float64

// Point rapresent a geojson point geometry object
type Point struct {
	Type       string `json:"type"`
	Coordinate `json:"coordinates"`
}

func (p *Point) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&p)
	return err
}

// Polygon rapresent a geojson polygon geometry object
type Polygon struct {
	Type        string `json:"type"`
	Coordinates `json:"coordinates"`
}

func (p *Polygon) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&p)
	return err
}

// MultiPolygon rapresent a geojson mulitpolygon  geometry object
type MultiPolygon struct {
	Type        string        `json:"type"`
	Coordinates []Coordinates `json:"coordinates"`
}

func (p *MultiPolygon) parseJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&p)
	return err
}

type geojson interface {
	parseJSON(io.Reader) error
}

// Endpoint is  the name of the geojson handler endpoint
const Endpoint = "/tos2/geojson/"

//Handler handles a request for a geojsonPoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// request
	resp, err := matcher(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// response
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
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
	case "polygon":
		p := Polygon{}
		err = p.parseJSON(r.Body)
		return &p, err
	case "point":
		p := Point{}
		err = p.parseJSON(r.Body)
		return &p, err
	case "multipolygon":
		p := MultiPolygon{}
		err = p.parseJSON(r.Body)
		return &p, err
	default:
		err = fmt.Errorf("Bad geoJSON object type")
	}
	return
}

func main() {
	// geojson
	geoJSONHandler := http.HandlerFunc(Handler)
	http.Handle(Endpoint, geoJSONHandler)
	// server
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)
	http.ListenAndServe(addr, nil)
}