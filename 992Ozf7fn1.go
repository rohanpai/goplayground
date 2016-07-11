package main

import (
	"fmt"
	"encoding/xml"
	)

type Gpx struct {
    Creator string `xml:"creator,attr"`
    Time string `xml:"metadata>time"`
    Title string `xml:"trk>name"`
    TrackPoints []TrackPoint `xml:"trk>trkseg>trkpt"`
}

type TrackPoint struct {
    Lat float64 `xml:"lat,attr"`
    Lon float64 `xml:"lon,attr"`
    Elevation float32 `xml:"ele"`
    Time string `xml:"time"`
    Temperature int `xml:"extensions>TrackPointExtension>atemp"`
}

func main() {
	data := `<gpx creator="StravaGPX" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" 

xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
 <metadata>
  <time>2013-02-16T10:11:25Z</time>
 </metadata>
 <trk>
  <name>Demo Data</name>
  <trkseg>
   <trkpt lat="51.6395658" lon="-3.3623858">
    <ele>111.6</ele>
    <time>2013-02-16T10:11:25Z</time>
    <extensions>
     <gpxtpx:TrackPointExtension>
      <gpxtpx:atemp>8</gpxtpx:atemp>
      <gpxtpx:hr>136</gpxtpx:hr>
      <gpxtpx:cad>0</gpxtpx:cad>
     </gpxtpx:TrackPointExtension>
    </extensions>
   </trkpt>
  </trkseg>
 </trk>
`
	g := &Gpx{}
	_ = xml.Unmarshal([]byte(data), g)
	fmt.Printf("len: %d\n", len(g.TrackPoints))
	fmt.Printf("temp: %v\n", g.TrackPoints[0].Temperature)
}