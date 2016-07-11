package main

import &#34;fmt&#34;
import &#34;math&#34;

const kmtomiles = float64(0.621371192)
const earthRadius = float64(6371)

func main() {
	var locationName [2]string
	var location [2][2]float64
	// York - lat,lon
	locationName[0] = &#34;York&#34;
	location[0][0] = 1.0803
	location[0][1] = 53.9583
	// Bristol - lat,lon
	locationName[1] = &#34;Bristol&#34;
	location[1][0] = 2.5833
	location[1][1] = 51.4500
	
	// Use haversine to get the resulting diatance between the two values
	var distance = Haversine(location[0][0], location[0][1], location[1][0], location[1][1])
	// We wish to use miles so will alter the resulting distance
	var distancemiles = distance * kmtomiles
	
	fmt.Printf(&#34;The distance between %s and %s is %.02f miles as the crow flies&#34;, locationName[0], locationName[1], distancemiles)
}

/*
 * The haversine formula will calculate the spherical distance as the crow flies 
 * between lat and lon for two given points in km
 */
func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)
	
	var a = math.Sin(deltaLat / 2) * math.Sin(deltaLat / 2) &#43; 
		math.Cos(latFrom * (math.Pi / 180)) * math.Cos(latTo * (math.Pi / 180)) *
		math.Sin(deltaLon / 2) * math.Sin(deltaLon / 2)
	var c = 2 * math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	
	distance = earthRadius * c
	
	return
}