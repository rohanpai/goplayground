package main

import (
	&#34;fmt&#34;
	&#34;encoding/xml&#34;
	)

type Gpx struct {
    Creator string `xml:&#34;creator,attr&#34;`
    Time string `xml:&#34;metadata&gt;time&#34;`
    Title string `xml:&#34;trk&gt;name&#34;`
    TrackPoints []TrackPoint `xml:&#34;trk&gt;trkseg&gt;trkpt&#34;`
}

type TrackPoint struct {
    Lat float64 `xml:&#34;lat,attr&#34;`
    Lon float64 `xml:&#34;lon,attr&#34;`
    Elevation float32 `xml:&#34;ele&#34;`
    Time string `xml:&#34;time&#34;`
    Temperature int `xml:&#34;extensions&gt;TrackPointExtension&gt;atemp&#34;`
}

func main() {
	data := `&lt;gpx creator=&#34;StravaGPX&#34; version=&#34;1.1&#34; xmlns=&#34;http://www.topografix.com/GPX/1/1&#34; 

xmlns:xsi=&#34;http://www.w3.org/2001/XMLSchema-instance&#34; xsi:schemaLocation=&#34;http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd&#34; xmlns:gpxtpx=&#34;http://www.garmin.com/xmlschemas/TrackPointExtension/v1&#34; xmlns:gpxx=&#34;http://www.garmin.com/xmlschemas/GpxExtensions/v3&#34;&gt;
 &lt;metadata&gt;
  &lt;time&gt;2013-02-16T10:11:25Z&lt;/time&gt;
 &lt;/metadata&gt;
 &lt;trk&gt;
  &lt;name&gt;Demo Data&lt;/name&gt;
  &lt;trkseg&gt;
   &lt;trkpt lat=&#34;51.6395658&#34; lon=&#34;-3.3623858&#34;&gt;
    &lt;ele&gt;111.6&lt;/ele&gt;
    &lt;time&gt;2013-02-16T10:11:25Z&lt;/time&gt;
    &lt;extensions&gt;
     &lt;gpxtpx:TrackPointExtension&gt;
      &lt;gpxtpx:atemp&gt;8&lt;/gpxtpx:atemp&gt;
      &lt;gpxtpx:hr&gt;136&lt;/gpxtpx:hr&gt;
      &lt;gpxtpx:cad&gt;0&lt;/gpxtpx:cad&gt;
     &lt;/gpxtpx:TrackPointExtension&gt;
    &lt;/extensions&gt;
   &lt;/trkpt&gt;
  &lt;/trkseg&gt;
 &lt;/trk&gt;
`
	g := &amp;Gpx{}
	_ = xml.Unmarshal([]byte(data), g)
	fmt.Printf(&#34;len: %d\n&#34;, len(g.TrackPoints))
	fmt.Printf(&#34;temp: %v\n&#34;, g.TrackPoints[0].Temperature)
}