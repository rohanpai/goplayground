package main

// Description:
//
// The program repeatedly downloads a given chunk of data
//     with a given interval in between.
// The time + the speed in Mbit/sec of each download are logged
//     in a CSV file having the current date in its name.
// The CSV file can be read later into Excel and a graph can be
//     made to show the variations in internet speed.

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func errExit(msg string) {
	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Println(`
	--------------------   SpeedTest foutmelding   ------------------
Speedtest: ` + msg + `

--------------------   SpeedTest help info   --------------------

Roep SpeedTest als volgt aan:
    speedtest <sample-interval> <sample-grootte>

Als u geen parameters opgeeft is dat hetzelfde als:
    speedtest 30 10

Het sample-interval mag niet kleiner zijn dan 30 seconden.

Kies een sample-grootte uit deze lijst:
    10, 25, 50, 100, 250, 500 of 1000

Het programma is nu gestopt!

-----------------------------------------------------------------

`)

	os.Exit(0)
}

func chkInterval(arg string) int {
	interval, err := strconv.Atoi(arg)
	if err != nil {
		errExit("Parameter 1: onjuiste waarde voor Sample Interval.")
	}
	if interval < 30 {
		errExit("Parameter 1: Sample Interval is te klein - moet >= 30 zijn.")
	}
	return interval
}

func chkSize(arg string) int {
	size, err := strconv.Atoi(arg)
	if err != nil {
		errExit("Parameter 2: onjuiste waarde voor Sample Grootte opgegeven.")
	}
	validSizes := [...]int{10, 25, 50, 100, 250, 500, 1000}
	found := false
	for _, value := range validSizes {
		if size == value {
			found = true
			break
		}
	}
	if found == false {
		errExit("Parameter 2: Sample grootte is ongeldig.")
	}
	return size
}

func main() {

	// Check number of args
	if len(os.Args) > 3 {
		errExit("Teveel input argumenten opgegeven.")
	}

	// Get Sample Interval in seconds from args
	interval := 30
	if len(os.Args) > 1 {
		interval = chkInterval(os.Args[1])
	}

	// Get Sample Size in MB from args
	size := 10
	if len(os.Args) > 2 {
		size = chkSize(os.Args[2])
	}

	// Setup output directory if it does not yet exist
	err := os.MkdirAll("c:\\speedtest", 0777)
	if err != nil {
		errExit("Fout bij aanmaken map 'c:\\speedtest'.")
	}

	// Setup some variables
	url := "http://speedtest.tweak.nl/" + strconv.Itoa(size) + "mb.bin"
	mbitsize := float64(size * 8)

	// Inform user
	dt := time.Now().Format(time.RFC3339)
	curdate := dt[:10]
	curtime := dt[11:19]
	fmt.Println("SpeedTest is gestart op " + curdate + " om " + curtime)
	fmt.Println("CSV bestand wordt opgeslagen in de map 'c:\\speedtest'")
	fmt.Println("Samplegrootte is " + strconv.Itoa(size) + " MB")
	fmt.Println("Sample interval is " + strconv.Itoa(interval) + " seconden")

	// Loop forever
	firstTime := true
	for {

		// Sleep during Sample Interval
		if !firstTime {
			time.Sleep(time.Duration(interval) * time.Second)
		}
		firstTime = false

		// Start timing the download
		start := time.Now()

		// Get the testfile
		resp, err := http.Get(url)
		if err != nil {
			errExit("Fout bij downloaden - http/get")
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			errExit("Probleem bij het downloaden - ioutil/readall")
		}

		// Calculate the elapsed time in seconds
		elapsed := time.Since(start)
		seconds := float64(elapsed) / float64(time.Second)

		// Internet speed in Mbit/Sec
		speednum := mbitsize / seconds
		speedstr := strconv.FormatFloat(speednum, 'f', 8, 64)

		// Prepare write to csv file
		dt = time.Now().Format(time.RFC3339)
		curdate = dt[:10]
		curtime = dt[11:19]
		csvfile := "c:\\speedtest\\speedtest_" + curdate + ".csv"

		// Open csv file in Append mode
		f, err := os.OpenFile(csvfile, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			fmt.Println(err)
			errExit("Fout bij openen csv bestand.")
		}

		// Format and show CSV record
		rec := curtime + ";" + strings.Replace(speedstr, ".", ",", 1)
		fmt.Println("Speedtest: " + rec)

		// Write CSV record
		bytes, err := f.WriteString(rec + "\n")
		if err != nil || bytes != (len(rec)+1) {
			errExit("Fout bij schrijven CSV bestand.")
		}

		// Close CSV file
		err = f.Close()
		if err != nil {
			errExit("Fout bij sluiten van het CSV bestand.")
		}

	}

	// Program can only be killed by the user (on purpose)
}
