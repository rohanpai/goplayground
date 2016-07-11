package main

// Description:
//
// The program repeatedly downloads a given chunk of data
//     with a given interval in between.
// The time &#43; the speed in Mbit/sec of each download are logged
//     in a CSV file having the current date in its name.
// The CSV file can be read later into Excel and a graph can be
//     made to show the variations in internet speed.

import (
	&#34;fmt&#34;
	&#34;io/ioutil&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
	&#34;time&#34;
)

func errExit(msg string) {
	fmt.Println(&#34;\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n&#34;)
	fmt.Println(`
	--------------------   SpeedTest foutmelding   ------------------
Speedtest: ` &#43; msg &#43; `

--------------------   SpeedTest help info   --------------------

Roep SpeedTest als volgt aan:
    speedtest &lt;sample-interval&gt; &lt;sample-grootte&gt;

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
		errExit(&#34;Parameter 1: onjuiste waarde voor Sample Interval.&#34;)
	}
	if interval &lt; 30 {
		errExit(&#34;Parameter 1: Sample Interval is te klein - moet &gt;= 30 zijn.&#34;)
	}
	return interval
}

func chkSize(arg string) int {
	size, err := strconv.Atoi(arg)
	if err != nil {
		errExit(&#34;Parameter 2: onjuiste waarde voor Sample Grootte opgegeven.&#34;)
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
		errExit(&#34;Parameter 2: Sample grootte is ongeldig.&#34;)
	}
	return size
}

func main() {

	// Check number of args
	if len(os.Args) &gt; 3 {
		errExit(&#34;Teveel input argumenten opgegeven.&#34;)
	}

	// Get Sample Interval in seconds from args
	interval := 30
	if len(os.Args) &gt; 1 {
		interval = chkInterval(os.Args[1])
	}

	// Get Sample Size in MB from args
	size := 10
	if len(os.Args) &gt; 2 {
		size = chkSize(os.Args[2])
	}

	// Setup output directory if it does not yet exist
	err := os.MkdirAll(&#34;c:\\speedtest&#34;, 0777)
	if err != nil {
		errExit(&#34;Fout bij aanmaken map &#39;c:\\speedtest&#39;.&#34;)
	}

	// Setup some variables
	url := &#34;http://speedtest.tweak.nl/&#34; &#43; strconv.Itoa(size) &#43; &#34;mb.bin&#34;
	mbitsize := float64(size * 8)

	// Inform user
	dt := time.Now().Format(time.RFC3339)
	curdate := dt[:10]
	curtime := dt[11:19]
	fmt.Println(&#34;SpeedTest is gestart op &#34; &#43; curdate &#43; &#34; om &#34; &#43; curtime)
	fmt.Println(&#34;CSV bestand wordt opgeslagen in de map &#39;c:\\speedtest&#39;&#34;)
	fmt.Println(&#34;Samplegrootte is &#34; &#43; strconv.Itoa(size) &#43; &#34; MB&#34;)
	fmt.Println(&#34;Sample interval is &#34; &#43; strconv.Itoa(interval) &#43; &#34; seconden&#34;)

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
			errExit(&#34;Fout bij downloaden - http/get&#34;)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			errExit(&#34;Probleem bij het downloaden - ioutil/readall&#34;)
		}

		// Calculate the elapsed time in seconds
		elapsed := time.Since(start)
		seconds := float64(elapsed) / float64(time.Second)

		// Internet speed in Mbit/Sec
		speednum := mbitsize / seconds
		speedstr := strconv.FormatFloat(speednum, &#39;f&#39;, 8, 64)

		// Prepare write to csv file
		dt = time.Now().Format(time.RFC3339)
		curdate = dt[:10]
		curtime = dt[11:19]
		csvfile := &#34;c:\\speedtest\\speedtest_&#34; &#43; curdate &#43; &#34;.csv&#34;

		// Open csv file in Append mode
		f, err := os.OpenFile(csvfile, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			fmt.Println(err)
			errExit(&#34;Fout bij openen csv bestand.&#34;)
		}

		// Format and show CSV record
		rec := curtime &#43; &#34;;&#34; &#43; strings.Replace(speedstr, &#34;.&#34;, &#34;,&#34;, 1)
		fmt.Println(&#34;Speedtest: &#34; &#43; rec)

		// Write CSV record
		bytes, err := f.WriteString(rec &#43; &#34;\n&#34;)
		if err != nil || bytes != (len(rec)&#43;1) {
			errExit(&#34;Fout bij schrijven CSV bestand.&#34;)
		}

		// Close CSV file
		err = f.Close()
		if err != nil {
			errExit(&#34;Fout bij sluiten van het CSV bestand.&#34;)
		}

	}

	// Program can only be killed by the user (on purpose)
}
