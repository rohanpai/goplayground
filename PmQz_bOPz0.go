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
//		GaryBoone's "GoStats" Statistics Package from https://github.com/GaryBoone/GoStats
//		csvdata/folder/GeographicalLocationXX.csv 
//		csvdata/folder/GeoTemporalMeasurementXX.csv
//			(where "folder" = pollutant type (Ozone, PM2.5)
//			(where "XX" represents the region desired; 01-10 and 25. 
//			       11 per pollutant)
//
//		Pollutant Types: "Ozone", "PM2.5"
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
//		        25 Mexico & Canada
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
//	References:	1. The constants, the function "AQMcompare", and everything below that 
//				function was written by Chris Tessum PhDc and modified 
//				by Steven Roste.
//			2. The "GoStats" Package was written by Gary Boone and is available at:
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
	"os"
	"log"
	"fmt"
	"encoding/csv"
	"strconv"
	"bitbucket.org/ctessum/gonetcdf"
	"errors"
	"math"
	"../GoStats"
	"time"
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
	dateFormat = "2006-01-02"
)

func main() {

	// Which Pollutant to evaluate
	var polntyp string
	fmt.Println("Enter the Pollution Type (Ozone or PM2.5): ")
	fmt.Scan(&polntyp)

	// Coordinates
	var SD string  // Start Date in YYYY-MM-DD
	var ED string  // End Date in YYYY-MM-DD
	fmt.Println("Enter the Start Date (YYYY-MM-DD): ")
	fmt.Scan(&SD)
	fmt.Println("Enter the Longitude Coordinate (YYYY-MM-DD): ")
	fmt.Scan(&ED)

	// Aquire Data
	fmt.Println("Aquiring measured EPA data...")
	valuDAT, locs := CSVreader(SD, ED, polntyp)
	fmt.Println("Aquiring modeled data...")
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
	fmt.Println("                                            ")
	fmt.Println("     Pollutant      Start Date     End Date")
	fmt.Println("     ", polntyp, "      ", SD, "  ", ED)
	fmt.Println("                                            ")
	fmt.Println("            Measured                   Modeled")
	fmt.Println("Mean:     ", meanDAT, "       ", meanMOD)
	fmt.Println("Min:      ", minDAT, "                   ", minMOD)
	fmt.Println("Max:      ", maxDAT, "                   ", maxMOD)
	fmt.Println("Count:    ", countDAT, "                        ",countMOD)
	fmt.Println("StdDev:   ", stdDAT, "        ", stdMOD)
	fmt.Println("Variance: ", varDAT, "        ", varMOD)
	fmt.Println("Skew:     ", skewDAT, "        ", skewMOD)
	fmt.Println("                                            ")
	fmt.Println("	  Error and Bias Statistics:")
	fmt.Println("Error:    ", e)
	fmt.Println("Abs Err:  ", eabs)
	fmt.Println("Bias:     ", bias)
	fmt.Println("Abs Bias: ", biasabs)
	fmt.Println("                                            ")

	Results(data) // "data" must be in [][]string format with all headers and data contained
}

fun Results(results [][]string) {

	// The results are output to a file that is named with a timestamp
	// so that if multiple files are created, the previous results are
	// not overwritten.

	var fn string
	fn = "Results"
	datetime := time.Now()
	filename := fmt.Sprint(fn, datetime.String(), ".csv") 

	fmt.Println("Creating File: ", filename, "...")
	file, err := os.Create(filename)
	if err != nil { 
		panic(err) 
	}
	defer file.Close()
	f := csv.NewWriter(file)
	fmt.Println("Writing Results to file...")
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
	for i := 1; i < 11; i++ {
	    filenum := i

	    ////////////////////////////////////////////////////////////////////////////////////
	    //
	    //	Determining which files to open
	    //
	    ////////////////////////////////////////////////////////////////////////////////////

	    var fileloc string
	    var filedat string

	    switch {
	    case polntyp == "Ozone":
		    switch {
			    case filenum == 1:
				 fileloc = "csvdata/Ozone/GeographicalLocation01.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement01.csv"
			    case filenum == 2:
				 fileloc = "csvdata/Ozone/GeographicalLocation02.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement02.csv"
			    case filenum == 3:
				 fileloc = "csvdata/Ozone/GeographicalLocation03.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement03.csv"
			    case filenum == 4:
				 fileloc = "csvdata/Ozone/GeographicalLocation04.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement04.csv"
			    case filenum == 5:
				 fileloc = "csvdata/Ozone/GeographicalLocation05.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement05.csv"
			    case filenum == 6:
				 fileloc = "csvdata/Ozone/GeographicalLocation06.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement06.csv"
			    case filenum == 7:
				 fileloc = "csvdata/Ozone/GeographicalLocation07.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement07.csv"
			    case filenum == 8:
				 fileloc = "csvdata/Ozone/GeographicalLocation08.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement08.csv"
			    case filenum == 9:
				 fileloc = "csvdata/Ozone/GeographicalLocation09.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement09.csv"
			    case filenum == 10:
				 fileloc = "csvdata/Ozone/GeographicalLocation10.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement10.csv"
			    case filenum == 11:
				 fileloc = "csvdata/Ozone/GeographicalLocation25.csv"
				 filedat = "csvdata/Ozone/GeoTemporalMeasurement25.csv"
			    default: 
				 panic("No file found.")
		    }
	    case polntyp == "PM2.5":
		    switch {
			    case filenum == 1:
				 fileloc = "csvdata/PM2.5/GeographicalLocation01.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement01.csv"
			    case filenum == 2:
				 fileloc = "csvdata/PM2.5/GeographicalLocation02.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement02.csv"
			    case filenum == 3:
				 fileloc = "csvdata/PM2.5/GeographicalLocation03.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement03.csv"
			    case filenum == 4:
				 fileloc = "csvdata/PM2.5/GeographicalLocation04.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement04.csv"
			    case filenum == 5:
				 fileloc = "csvdata/PM2.5/GeographicalLocation05.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement05.csv"
			    case filenum == 6:
				 fileloc = "csvdata/PM2.5/GeographicalLocation06.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement06.csv"
			    case filenum == 7:
				 fileloc = "csvdata/PM2.5/GeographicalLocation07.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement07.csv"
			    case filenum == 8:
				 fileloc = "csvdata/PM2.5/GeographicalLocation08.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement08.csv"
			    case filenum == 9:
				 fileloc = "csvdata/PM2.5/GeographicalLocation09.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement09.csv"
			    case filenum == 10:
				 fileloc = "csvdata/PM2.5/GeographicalLocation10.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement10.csv"
			    case filenum == 11:
				 fileloc = "csvdata/PM2.5/GeographicalLocation25.csv"
				 filedat = "csvdata/PM2.5/GeoTemporalMeasurement25.csv"
			    default: 
				 panic("No file found.")
		    }
	    default:
		 panic("Pollutant data not found.")
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
		    log.Fatalf("error reading all lines: %v", err)
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
		fmt.Println("[ID Tag]", label1, label2, label3, label4, label5, label6, label7, label8, label9)
		for w := 0; w < len(idtag)-2; w++ { // loop through all data
			if SDt <= time.Parse(dateFormat, date[w]) && EDt >= time.Parse(dateFormat, date[w]) {
				fmt.Println(" ", idtag[w], "     ", lats[w], "         ", lngs[w], "      ", poln[w], " ", date[w], " ", time[w], " ", zone[w], " ", valu[w], " ", unit[w], " ", resl[w]) // print to screen data table
				if count == 0 {
					valuOUT[count] = valu[w]
				} else {
					valuOUT = append(valuOUT, valu[w])
				}
				count = count + 1
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
	if polntyp == "PM2.5" {
		flag = 1
	}
	count := 0
	valuOUT :=make([]float64,1)
	var SDt time.Time
	var EDt time.Time

	pols := [...]string{"PM2_5_DRY", "o3"}

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

		fn := "wrfout/wrfout_d01_"
		cdate := i.Format(dateFormat)
		filename := fmt.Sprint(fn, cdate, "_00_00_00")
	
		// Open WRF file, handle any errors
		f, e := gonetcdf.Open(filename, "nowrite")
		if e != nil {
			fmt.Println("No wrf file available for: ", cdate)
			continue
		}

		// get the necessary data for calculating layer heights
		fmt.Println("Getting the base state geopotential...")
		PHB, e := NewWRFarray(f, "PHB")
		if e != nil {
			panic(e)
		}
		fmt.Println("Getting the perturbation geopotential...")
		PH, e := NewWRFarray(f, "PH")
		if e != nil {
			panic(e)
		}

		for _, pol := range pols {
			fmt.Printf("Getting data for %v...\n", pol)
			data, e := NewWRFarray(f, pol)
			if e != nil {
				panic(e)
			}
			for i := 0; i < len(locs)-1; i++ { 		
				lat, lon := locs(i)
				for h := 0; h < 24; h++ {
					height := 50. // meters

					/////////////////////////////////////////////////////

					ycell, xcell, e := Lambert(lat, lon, reflat, reflon, phi1,
						phi2, x_o, y_o, R, dx, dy, nx, ny)
					if e != nil {
						panic(e)
					}

					layer, e := LayerHeight(height, ycell, xcell, h, &PH, &PHB)

					index := []int{h, layer, ycell, xcell}
					conc := data.GetVal(index)

					// Which values to store
					if flag == 0 && pol == "o3" {
						// store Ozone values
						if count == 0 {
							valuOUT[h] = conc
						} else {
							valuOUT = append(valuOUT, conc)
						}
						// print the result
						fmt.Printf("The value for pollutant %v at %v, %v for hour %v in layer %v is %v\n",
						pol, lon, lat, h, layer, conc)
						count = count + 1
					} else if flag == 1 && pol == "PM2_5_DRY" {
						//store PM2.5 values
						if count == 0 {
							valuOUT[h] = conc
						} else {
							valuOUT = append(valuOUT, conc)
						}
						// print the result
						fmt.Printf("The value for pollutant %v at %v, %v for hour %v in layer %v is %v\n",
						pol, lon, lat, h, layer, conc)
						count = count + 1
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
	for i := 0; i < len(index); i++ {
		mul := 1
		for j := i + 1; j < len(index); j++ {
			mul = mul * d.Dims[j]
		}
		index1d = index1d + index[i]*mul
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
		math.Log(math.Tan(0.25*math.Pi+0.5*phi2)/
			math.Tan(0.25*math.Pi+0.5*phi1))

	F := math.Cos(phi1) * math.Pow(
		math.Tan(0.25*math.Pi+0.5*phi1), n) / n

	rho := F * math.Pow((1./math.Tan(0.25*math.Pi+0.5*lat)), n)

	rho_o := F * math.Pow(1./math.Tan(0.25*math.Pi+0.5*reflat), n)

	x := rho * math.Sin(n*(lon-reflon)) * R
	y := (rho_o - rho*math.Cos(n*(lon-reflon))) * R

	xcell = int(math.Floor((x - x_o) / dx))
	ycell = int(math.Floor((y - y_o) / dy))

	if xcell >= nx || ycell >= ny || xcell < 0 || ycell < 0 {
		fmt.Println(xcell)
		fmt.Println(ycell)
		err = errors.New("Cell number out of bounds")
	}

	return
}

func LayerHeight(height float64, ycell int, xcell int, hour int,
	PH *WRFarray, PHB *WRFarray) (layer int, err error) {
	for k := 1; k < PH.Dims[1]; k++ {
		PHi := []int{hour, k, ycell, xcell}
		PHBi := []int{hour, k , ycell, xcell}
		layerHeight := (PH.GetVal(PHi) + PHB.GetVal(PHBi)) / g
		if layerHeight < height {
			layer = k - 1
		} else {
			return
		}
	}
	errors.New("Height is above top layer")
	return
}