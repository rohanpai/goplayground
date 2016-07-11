//////////////////////////////////////////////////////////////////////////////////////////////
//
//	Usage:		go run AQMstat.go
//
//	Summary:	This program was written for use to statistically compare modeled and 
//		measured values of ambient air pollutants.
//		Data time range: Jan-01-2005 through Dec-31-2005
//			ISO 8601 Format: 2005-01-01 through 2005-12-31
//
//	Inputs:
//		SD = prompts user for the start date as a string in ISO 8601 format YYYY-MM-DD
//		ED = prompts user for the end date as a string in ISO 8601 format YYYY-MM-D
//		polntyp = prompts user for the type of pollution to be evaluated
//
//	Outputs:
//		ResultsYYYY-MM-DD HH:MM:SS.ssss 0000 TIMEZONE.csv = a CSV file of the results 
//			is created with a filename that includes a time-stamp of when the 
//			results were created with the date-time-timezone format shown above. 
//
//	Associated Files:
//		GoNetCDF Package from http://www.unidata.ucar.edu/software/netcdf/docs/netcdf-c/
//		GaryBoone&#39;s &#34;GoStats&#34; Statistics Package from https://github.com/GaryBoone/GoStats
//		csvdata/folder/GeographicalLocationXX.csv 
//		csvdata/folder/GeoTemporalMeasurementXX.csv
//			(where &#34;folder&#34; = pollutant type (Ozone, PM2.5)
//			(where &#34;XX&#34; represents the region desired; 01-10 and 25. 
//			       11 per pollutant)
//
//		Pollutant Types: &#34;Ozone&#34;, &#34;PM2.5&#34;
//
//		Region Identifiers:
//		        01 Boston
//		        02 New York City
//		        03 Philidelphia
//		        04 Atlanta
//		        05 Chicago
//		        06 Dallas
//		        07 Kansas City
//		        08 Denver
//		        09 San Fransisco
//		        10 Seattle
//		        25 Mexico &amp; Canada
//
//	Version:	2012.06.25
//
//	Written by:	Steven J. Roste, EIT - June 2012
//			Junior Scientist 
//			Department of Civil Engineering
//			College of Science and Engineering
//			University of Minnesota
//			roste025@umn.edu
//
//	References:	1. The constants, the function &#34;AQMcompare&#34;, and everything below that 
//				function was written by Chris Tessum PhDc and modified 
//				by Steven Roste.
//			2. The &#34;GoStats&#34; Package was written by Gary Boone and is available at:
//				https://github.com/GaryBoone/GoStats
//			3. GoNetCDF Package from 
//				http://www.unidata.ucar.edu/software/netcdf/docs/netcdf-c/
//			4. http://golang.org/
//			5. http://stackoverflow.com/
//			6. https://groups.google.com/forum/?fromgroups#!forum/golang-nuts
//			7. Summerfield, Mark. Programming in Go: Creating Applications for the 21st Century. 
//				Upper Saddle River, NJ: Addison-Wesley, 2012. Print.
//			8. Chisnall, David. The Go Programming Language: Phrasebook. 
//				Upper Saddle River, NJ: Addison-Wesley Professional, 2012. Print.
//
//////////////////////////////////////////////////////////////////////////////////////////////

package main

import(
	&#34;os&#34;
	&#34;log&#34;
	&#34;fmt&#34;
	&#34;encoding/csv&#34;
	&#34;strconv&#34;
	&#34;bitbucket.org/ctessum/gonetcdf&#34;
	&#34;errors&#34;
	&#34;math&#34;
	&#34;../GoStats&#34;
	&#34;time&#34;
)

//  some information about the projection
const (
	reflat = 40.0
	reflon = -97.0
	phi1   = 45.0
	phi2   = 33.0
	x_o    = -2736000.
	y_o    = -2088000.
	R      = 6370997.
	dx     = 12000.
	dy     = 12000.
	nx     = 444
	ny     = 336
	g      = 9.8 // m/s2
	dateFormat = &#34;2006-01-02&#34;
)

func main() {

	// Which Pollutant to evaluate
	var polntyp string
	fmt.Println(&#34;Enter the Pollution Type (Ozone or PM2.5): &#34;)
	fmt.Scan(&amp;polntyp)

	// Coordinates
	var SD string  // Start Date in YYYY-MM-DD
	var ED string  // End Date in YYYY-MM-DD
	fmt.Println(&#34;Enter the Start Date (YYYY-MM-DD): &#34;)
	fmt.Scan(&amp;SD)
	fmt.Println(&#34;Enter the Longitude Coordinate (YYYY-MM-DD): &#34;)
	fmt.Scan(&amp;ED)

	// Aquire Data
	fmt.Println(&#34;Aquiring measured EPA data...&#34;)
	valuDAT, locs := CSVreader(SD, ED, polntyp)
	fmt.Println(&#34;Aquiring modeled data...&#34;)
	valuMOD := AQMcompare(SD, ED, polntyp, locs)

	// Process measured data
	meanDAT := stats.StatsMean(valuDAT)
	minDAT := stats.StatsMin(valuDAT)
	maxDAT := stats.StatsMax(valuDAT)
	countDAT := stats.StatsCount(valuDAT)
	stdDAT := stats.StatsSampleStandardDeviation(valuDAT)
	varDAT := stats.StatsSampleVariance(valuDAT)
	skewDAT := stats.StatsSampleSkew(valuDAT)

	// Process modeled data
	meanMOD := stats.StatsMean(valuMOD)
	minMOD := stats.StatsMin(valuMOD)
	maxMOD := stats.StatsMax(valuMOD)
	countMOD := stats.StatsCount(valuMOD)
	stdMOD := stats.StatsSampleStandardDeviation(valuMOD)
	varMOD := stats.StatsSampleVariance(valuMOD)
	skewMOD := stats.StatsSampleSkew(valuMOD)

	// Comparitive Stats
	e := meanMOD - meanDAT
	eabs := math.Abs(e)
	bias := e / meanDAT
	biasabs := math.Abs(bias)
	

	// Print Stats to screen
	fmt.Println(&#34;                                            &#34;)
	fmt.Println(&#34;     Pollutant      Start Date     End Date&#34;)
	fmt.Println(&#34;     &#34;, polntyp, &#34;      &#34;, SD, &#34;  &#34;, ED)
	fmt.Println(&#34;                                            &#34;)
	fmt.Println(&#34;            Measured                   Modeled&#34;)
	fmt.Println(&#34;Mean:     &#34;, meanDAT, &#34;       &#34;, meanMOD)
	fmt.Println(&#34;Min:      &#34;, minDAT, &#34;                   &#34;, minMOD)
	fmt.Println(&#34;Max:      &#34;, maxDAT, &#34;                   &#34;, maxMOD)
	fmt.Println(&#34;Count:    &#34;, countDAT, &#34;                        &#34;,countMOD)
	fmt.Println(&#34;StdDev:   &#34;, stdDAT, &#34;        &#34;, stdMOD)
	fmt.Println(&#34;Variance: &#34;, varDAT, &#34;        &#34;, varMOD)
	fmt.Println(&#34;Skew:     &#34;, skewDAT, &#34;        &#34;, skewMOD)
	fmt.Println(&#34;                                            &#34;)
	fmt.Println(&#34;	  Error and Bias Statistics:&#34;)
	fmt.Println(&#34;Error:    &#34;, e)
	fmt.Println(&#34;Abs Err:  &#34;, eabs)
	fmt.Println(&#34;Bias:     &#34;, bias)
	fmt.Println(&#34;Abs Bias: &#34;, biasabs)
	fmt.Println(&#34;                                            &#34;)

	Results(data) // &#34;data&#34; must be in [][]string format with all headers and data contained
}

fun Results(results [][]string) {

	// The results are output to a file that is named with a timestamp
	// so that if multiple files are created, the previous results are
	// not overwritten.

	var fn string
	fn = &#34;Results&#34;
	datetime := time.Now()
	filename := fmt.Sprint(fn, datetime.String(), &#34;.csv&#34;) 

	fmt.Println(&#34;Creating File: &#34;, filename, &#34;...&#34;)
	file, err := os.Create(filename)
	if err != nil { 
		panic(err) 
	}
	defer file.Close()
	f := csv.NewWriter(file)
	fmt.Println(&#34;Writing Results to file...&#34;)
	f.WriteAll(results)
}

func CSVreader(SD string, ED string, polntyp string) []float64, [][]float64 {

	count := 0
	valuOUT := make([]float64,1)
	var SDt time.Time
	var EDt time.Time
	SDt, err = time.Parse(dateFormat, SD)
	if err != nil {
		panic(err)
	}
	EDt, err = time.Parse(dateFormat, ED)
	if err != nil {
		panic(err)
	}
	EDf := EDt.AddDate(0,0,1)

	//Cycle through all Location and Data files
	for i := 1; i &lt; 11; i&#43;&#43; {
	    filenum := i

	    ////////////////////////////////////////////////////////////////////////////////////
	    //
	    //	Determining which files to open
	    //
	    ////////////////////////////////////////////////////////////////////////////////////

	    var fileloc string
	    var filedat string

	    switch {
	    case polntyp == &#34;Ozone&#34;:
		    switch {
			    case filenum == 1:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation01.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement01.csv&#34;
			    case filenum == 2:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation02.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement02.csv&#34;
			    case filenum == 3:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation03.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement03.csv&#34;
			    case filenum == 4:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation04.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement04.csv&#34;
			    case filenum == 5:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation05.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement05.csv&#34;
			    case filenum == 6:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation06.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement06.csv&#34;
			    case filenum == 7:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation07.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement07.csv&#34;
			    case filenum == 8:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation08.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement08.csv&#34;
			    case filenum == 9:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation09.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement09.csv&#34;
			    case filenum == 10:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation10.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement10.csv&#34;
			    case filenum == 11:
				 fileloc = &#34;csvdata/Ozone/GeographicalLocation25.csv&#34;
				 filedat = &#34;csvdata/Ozone/GeoTemporalMeasurement25.csv&#34;
			    default: 
				 panic(&#34;No file found.&#34;)
		    }
	    case polntyp == &#34;PM2.5&#34;:
		    switch {
			    case filenum == 1:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation01.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement01.csv&#34;
			    case filenum == 2:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation02.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement02.csv&#34;
			    case filenum == 3:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation03.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement03.csv&#34;
			    case filenum == 4:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation04.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement04.csv&#34;
			    case filenum == 5:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation05.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement05.csv&#34;
			    case filenum == 6:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation06.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement06.csv&#34;
			    case filenum == 7:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation07.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement07.csv&#34;
			    case filenum == 8:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation08.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement08.csv&#34;
			    case filenum == 9:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation09.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement09.csv&#34;
			    case filenum == 10:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation10.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement10.csv&#34;
			    case filenum == 11:
				 fileloc = &#34;csvdata/PM2.5/GeographicalLocation25.csv&#34;
				 filedat = &#34;csvdata/PM2.5/GeoTemporalMeasurement25.csv&#34;
			    default: 
				 panic(&#34;No file found.&#34;)
		    }
	    default:
		 panic(&#34;Pollutant data not found.&#34;)
	    }	

	    ////////////////////////////////////////////////////////////////////////////////////
	    //
	    //	Extracting Locations from the GeographicalLocaionXX.csv file
	    //
	    ////////////////////////////////////////////////////////////////////////////////////

	    // Open Location CSV file, handle any errors
	    filename1, err := os.Open(fileloc) // folder/pollutant/filename.csv
	    if err != nil { 
		    panic(err) 
	    }
	    defer filename1.Close() // close file when finished
	    r1 := csv.NewReader(filename1) 
	    lines1, err := r1.ReadAll() // Read to EOF and store in lines2
	    if err != nil {
		    panic(err)
	    }

	    // define data variables to be extracted
	    idtag := make([]int, len(lines1)-1) // define file Id tag
	    label1 := make([]string, 2)
	    label2 := make([]string, 2)
	    lats := make([]float64, len(lines1)-1) // define lats
	    lngs := make([]float64, len(lines1)-1) // define lngs
	    // loop through all rows to extract data and convert to float64 values
	    for j, row1 := range lines1 { 
		    if j == 0 {
			    label1[1] = row1[0]
			    label2[1] = row1[1]	
		    } else {
		       	    lats[j-1], err = strconv.ParseFloat(row1[0],64) // store row value from column 1
			    if err != nil {
		    		panic(err)
	    		    }
		       	    lngs[j-1], err = strconv.ParseFloat(row1[1],64) // store row value from column 2
			    if err != nil {
		    		panic(err)
	    		    }
		       	    idtag[j-1], err = strconv.Atoi(row1[3]) // store file Id tag number for reference
			    if err != nil {
		    		panic(err)
	    		    }
		    }
	    }

	    ////////////////////////////////////////////////////////////////////////////////////
	    //
	    //	Extracting Data from the corresponding GeoTemporalMeasurementXX.csv file
	    //
	    ////////////////////////////////////////////////////////////////////////////////////

	    // extract correlated data from same region

	    // Open Measurement CSV file, handle any errors
	    filename2, err := os.Open(filedat) // folder/pollutant/filename.csv
	    if err != nil { 
		    panic(err) 
	    }
	    defer filename2.Close() // close file when finished
	    r2 := csv.NewReader(filename2) 
	    lines2, err := r2.ReadAll() // Read to EOF and store in lines2
	    if err != nil {
		    log.Fatalf(&#34;error reading all lines: %v&#34;, err)
	    }

	    // define data variables to be extracted
	    label3 := make([]string, 2)
	    label4 := make([]string, 2)
	    label5 := make([]string, 2)
	    label6 := make([]string, 2)
	    label7 := make([]string, 2)
	    label8 := make([]string, 2)
	    label9 := make([]string, 2)
	    poln := make([]string, len(lines2)-1) // pollution type
	    date := make([]string, len(lines2)-1) // time of measurement
	    time := make([]string, len(lines2)-1) // local time
	    zone := make([]string, len(lines2)-1) // local time zone
	    valu := make([]float64, len(lines2)-1) // measurment value
	    unit := make([]string, len(lines2)-1) // unit of measurement
	    resl := make([]string, len(lines2)-1) // time resolution

	    // loop through all rows to extract data and convert measurements to float64 values
	    for k, line2 := range lines2 { 
		    if k == 0 {
			    label3[1] = line2[1]
			    label4[1] = line2[2]
			    label5[1] = line2[3]
			    label6[1] = line2[4]
			    label7[1] = line2[6]
			    label8[1] = line2[7]
			    label9[1] = line2[12]			
			    continue
		    } else {
			    poln[k-1] = line2[1] // store row values from column 2
			    date[k-1] = line2[2] // store row values from column 3
			    time[k-1] = line2[3] // etc
			    zone[k-1] = line2[4]
			    valu[k-1], err = strconv.ParseFloat(line2[6],64) // store row value from column 2
			    if err != nil {
		    		panic(err)
	    		    }
			    unit[k-1] = line2[7]
			    resl[k-1] = line2[12]
		    }
	    }
		////////////////////////////////////////////////////////////////////////////////////
		//
		//	Display relevant data in a table format to command window
		//
		////////////////////////////////////////////////////////////////////////////////////

		// display all extracted data
		fmt.Println(fileloc, filedat)
		fmt.Println(&#34;[ID Tag]&#34;, label1, label2, label3, label4, label5, label6, label7, label8, label9)
		for w := 0; w &lt; len(idtag)-2; w&#43;&#43; { // loop through all data
			if SDt &lt;= time.Parse(dateFormat, date[w]) &amp;&amp; EDt &gt;= time.Parse(dateFormat, date[w]) {
				fmt.Println(&#34; &#34;, idtag[w], &#34;     &#34;, lats[w], &#34;         &#34;, lngs[w], &#34;      &#34;, poln[w], &#34; &#34;, date[w], &#34; &#34;, time[w], &#34; &#34;, zone[w], &#34; &#34;, valu[w], &#34; &#34;, unit[w], &#34; &#34;, resl[w]) // print to screen data table
				if count == 0 {
					valuOUT[count] = valu[w]
				} else {
					valuOUT = append(valuOUT, valu[w])
				}
				count = count &#43; 1
			}
		}
	}

	return valuOUT, locs
}

/////////////////////////////////////////////////////////////////
//
//	AQMcompare and related functions
//
/////////////////////////////////////////////////////////////////

func AQMcompare(SD string, ED string, polntyp string, locs [][]float64) []float64 {

	flag := 0 // 0 for Ozone, 1 for PM2.5
	if polntyp == &#34;PM2.5&#34; {
		flag = 1
	}
	count := 0
	valuOUT :=make([]float64,1)
	var SDt time.Time
	var EDt time.Time

	pols := [...]string{&#34;PM2_5_DRY&#34;, &#34;o3&#34;}

	SDt, err = time.Parse(dateFormat, SD)
	if err != nil {
		panic(err)
	}
	EDt, err = time.Parse(dateFormat, ED)
	if err != nil {
		panic(err)
	}

	EDf := EDt.AddDate(0,0,1)

	for i := SDt; i != EDf;  {

		fn := &#34;wrfout/wrfout_d01_&#34;
		cdate := i.Format(dateFormat)
		filename := fmt.Sprint(fn, cdate, &#34;_00_00_00&#34;)
	
		// Open WRF file, handle any errors
		f, e := gonetcdf.Open(filename, &#34;nowrite&#34;)
		if e != nil {
			fmt.Println(&#34;No wrf file available for: &#34;, cdate)
			continue
		}

		// get the necessary data for calculating layer heights
		fmt.Println(&#34;Getting the base state geopotential...&#34;)
		PHB, e := NewWRFarray(f, &#34;PHB&#34;)
		if e != nil {
			panic(e)
		}
		fmt.Println(&#34;Getting the perturbation geopotential...&#34;)
		PH, e := NewWRFarray(f, &#34;PH&#34;)
		if e != nil {
			panic(e)
		}

		for _, pol := range pols {
			fmt.Printf(&#34;Getting data for %v...\n&#34;, pol)
			data, e := NewWRFarray(f, pol)
			if e != nil {
				panic(e)
			}
			for i := 0; i &lt; len(locs)-1; i&#43;&#43; { 		
				lat, lon := locs(i)
				for h := 0; h &lt; 24; h&#43;&#43; {
					height := 50. // meters

					/////////////////////////////////////////////////////

					ycell, xcell, e := Lambert(lat, lon, reflat, reflon, phi1,
						phi2, x_o, y_o, R, dx, dy, nx, ny)
					if e != nil {
						panic(e)
					}

					layer, e := LayerHeight(height, ycell, xcell, h, &amp;PH, &amp;PHB)

					index := []int{h, layer, ycell, xcell}
					conc := data.GetVal(index)

					// Which values to store
					if flag == 0 &amp;&amp; pol == &#34;o3&#34; {
						// store Ozone values
						if count == 0 {
							valuOUT[h] = conc
						} else {
							valuOUT = append(valuOUT, conc)
						}
						// print the result
						fmt.Printf(&#34;The value for pollutant %v at %v, %v for hour %v in layer %v is %v\n&#34;,
						pol, lon, lat, h, layer, conc)
						count = count &#43; 1
					} else if flag == 1 &amp;&amp; pol == &#34;PM2_5_DRY&#34; {
						//store PM2.5 values
						if count == 0 {
							valuOUT[h] = conc
						} else {
							valuOUT = append(valuOUT, conc)
						}
						// print the result
						fmt.Printf(&#34;The value for pollutant %v at %v, %v for hour %v in layer %v is %v\n&#34;,
						pol, lon, lat, h, layer, conc)
						count = count &#43; 1
					}
				}
			}
		}
		i = i.AddDate(0,0,1)
	}
	return valuOUT
}

func NewWRFarray(f gonetcdf.NCfile, name string) (
	out WRFarray, e error) {
	out.Dims, e = f.VarSize(name)
	if e != nil {
		return
	}
	out.Data, e = f.GetVarDouble(name)
	if e != nil {
		return
	}
	return
}

type WRFarray struct {
	Data []float64
	Dims []int
	Name string
}

// Function IndexTo1D takes an array of indecies for a
// multi-dimensional array and the dimensions of that array,
// calculates the 1D-array index, and returns the corresponding value.
func (d *WRFarray) GetVal(index []int) (val float64) {
	index1d := 0
	for i := 0; i &lt; len(index); i&#43;&#43; {
		mul := 1
		for j := i &#43; 1; j &lt; len(index); j&#43;&#43; {
			mul = mul * d.Dims[j]
		}
		index1d = index1d &#43; index[i]*mul
	}
	return d.Data[index1d]
}

// Function Lambert takes the latitude and longitude of a point plus
// the definition of a lambert conic conformal projection as inputs
// (reflat and reflon are the reference coordinates of the projection,
// phi1 and phi2 are the two other reference latitudes, x_o and y_o
// are the coordinates in meters of the lower left corner of the 
// grid (in meters), R is the radius of the Earth (in meters), 
// dx and dy specify the size of each grid cell (in meters), and 
// nx and ny specify the number of grid cells in each direction.
// The function returns returns the location of the grid cell
// containing the input latitude.
func Lambert(lat float64, lon float64, reflat float64, reflon float64,
	phi1 float64, phi2 float64, x_o float64, y_o float64,
	R float64, dx float64, dy float64, nx int, ny int) (
	ycell int, xcell int, err error) {

	lon = lon * math.Pi / 180.
	lat = lat * math.Pi / 180.
	reflat = reflat * math.Pi / 180.
	reflon = reflon * math.Pi / 180.
	phi1 = phi1 * math.Pi / 180.
	phi2 = phi2 * math.Pi / 180.

	n := math.Log(math.Cos(phi1)/math.Cos(phi2)) /
		math.Log(math.Tan(0.25*math.Pi&#43;0.5*phi2)/
			math.Tan(0.25*math.Pi&#43;0.5*phi1))

	F := math.Cos(phi1) * math.Pow(
		math.Tan(0.25*math.Pi&#43;0.5*phi1), n) / n

	rho := F * math.Pow((1./math.Tan(0.25*math.Pi&#43;0.5*lat)), n)

	rho_o := F * math.Pow(1./math.Tan(0.25*math.Pi&#43;0.5*reflat), n)

	x := rho * math.Sin(n*(lon-reflon)) * R
	y := (rho_o - rho*math.Cos(n*(lon-reflon))) * R

	xcell = int(math.Floor((x - x_o) / dx))
	ycell = int(math.Floor((y - y_o) / dy))

	if xcell &gt;= nx || ycell &gt;= ny || xcell &lt; 0 || ycell &lt; 0 {
		fmt.Println(xcell)
		fmt.Println(ycell)
		err = errors.New(&#34;Cell number out of bounds&#34;)
	}

	return
}

func LayerHeight(height float64, ycell int, xcell int, hour int,
	PH *WRFarray, PHB *WRFarray) (layer int, err error) {
	for k := 1; k &lt; PH.Dims[1]; k&#43;&#43; {
		PHi := []int{hour, k, ycell, xcell}
		PHBi := []int{hour, k , ycell, xcell}
		layerHeight := (PH.GetVal(PHi) &#43; PHB.GetVal(PHBi)) / g
		if layerHeight &lt; height {
			layer = k - 1
		} else {
			return
		}
	}
	errors.New(&#34;Height is above top layer&#34;)
	return
}