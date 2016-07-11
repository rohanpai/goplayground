package main

import (
	&#34;bufio&#34;
	&#34;fmt&#34;
	gzip &#34;github.com/klauspost/pgzip&#34;
	//&#34;compress/gzip&#34;
	&#34;os&#34;
	&#34;runtime&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
	&#34;time&#34;
)

// Time processing allcities.txt from http://download.geonames.org/export/dump/ - recompressed to gzip.
//
// pgzip: Processed 10013208 entries in 43.484375s, 230271.4 entries/sec.
// gzip: Processed 10013208 entries in 1m28.859375s, 112686.0 entries/sec
//
// Diff: 117585.4 = 104%
type City struct {
	GeonameID      int
	Name           string   //  name of geographical point (utf8) varchar(200)
	AsciiName      string   //  name of geographical point in plain ascii characters, varchar(200)
	AlternateNames []string //  alternatenames, comma separated, ascii names automatically transliterated, convenience attribute from alternatename table, varchar(10000)
	Latitude       float64  //  latitude in decimal degrees (wgs84)
	Longitude      float64  //  longitude in decimal degrees (wgs84)
	FeatureClass   string   //  see httpstring // //www.geonames.org/export/codes.html, char(1)
	FeatureCode    string   //  see httpstring // //www.geonames.org/export/codes.html, varchar(10)
	CountryCode    string   //  ISO-3166 2-letter country code, 2 characters
	CC2            string   //  alternate country codes, comma separated, ISO-3166 2-letter country code, 60 characters
	/* Skip 4 Admin */
	Population int       //  bigint (8 byte int)
	Elevation  int       //  in meters, integer
	Dem        string    //  digital elevation model, srtm3 or gtopo30, average elevation of 3&#39;&#39;x3&#39;&#39; (ca 90mx90m) or 30&#39;&#39;x30&#39;&#39; (ca 900mx900m) area in meters, integer. srtm processed by cgiar/ciat.
	Timezone   string    //  the timezone id (see file timeZone.txt) varchar(40)
	Modified   time.Time //  date of last modification in yyyy-MM-dd format
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.Open(&#34;allCountries.txt.gz&#34;)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(r)
	t := time.Now()

	n := 0
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, &#34;\t&#34;)
		if len(s) &lt; 19 {
			continue
		}
		c := City{}
		c.GeonameID, _ = strconv.Atoi(s[0])
		c.Name = s[1]
		c.AsciiName = s[2]
		c.AlternateNames = strings.Split(s[3], &#34;,&#34;)
		c.Latitude, _ = strconv.ParseFloat(s[4], 64)
		c.Longitude, _ = strconv.ParseFloat(s[5], 64)
		c.FeatureClass = s[6]
		c.FeatureCode = s[7]
		c.CountryCode = s[8]
		c.CC2 = s[9]
		c.Population, _ = strconv.Atoi(s[14])
		c.Elevation, _ = strconv.Atoi(s[15])
		c.Dem = s[16]
		c.Timezone = s[17]
		c.Modified, _ = time.Parse(&#34;2006-01-02&#34;, s[18])
		//fmt.Printf(&#34;%#v\n&#34;, c)
		n&#43;&#43;
	}
	d := time.Since(t)
	fmt.Printf(&#34;Processed %d entries in %v, %.1f entries/sec.&#34;, n, d, float64(n)/(float64(d)/float64(time.Second)))
}
