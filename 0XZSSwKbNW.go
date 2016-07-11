// This package uses a Levenshtein Distance implementation to compare customer accounts to help find duplicates

package main

import (
	&#34;container/ring&#34;
	&#34;encoding/csv&#34;
	&#34;errors&#34;
	&#34;files&#34;
	&#34;fmt&#34;
	&#34;io&#34;
	&#34;os&#34;
	&#34;runtime&#34;
	&#34;strconv&#34;
)

// a simple struct representing a customer account, by account number and name
type customer struct {
	num  string
	name string
}

// a job represents all the customers in a particular region
// each customer in the region will be compared to each other in the same region
type job struct {
	region   string
	accounts []customer
	comm     chan []string
}

// some consts for column numbers from the (csv) input
const (
	accountField = 0
	nameField    = 1
	regionField  = 5
)

// a constant for the output buffer in func main
const bufSize = 100

// function ld is the Levenshtein Distance implementation
// this was taken from Rosetta Code at: http://rosettacode.org/wiki/Levenshtein_distance#Go
func ld(s, t string) int {
	d := make([][]int, len(s)&#43;1)
	for i := range d {
		d[i] = make([]int, len(t)&#43;1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j &lt;= len(t); j&#43;&#43; {
		for i := 1; i &lt;= len(s); i&#43;&#43; {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] &lt; min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] &lt; min {
					min = d[i-1][j-1]
				}
				d[i][j] = min &#43; 1
			}
		}

	}
	return d[len(s)][len(t)]
}

// this function reads the customer data (in csv format) from a given io.Reader and returns them in a map
// the map uses the region code (a string) as the key and a slice of customers as the value
// this is to group the customers together by region, as we will not be comparing customers from different regions
func read(src io.Reader) map[string][]customer {
	// make a map to store the customer records, organized by region
	regions := make(map[string][]customer)

	// read all records and organize them into the regions map
	e := files.ReadCSV(src,
		func(line int, fields []string) error {
			if line &lt; 2 {
				return nil
			}
			if len(fields) != 25 {
				return errors.New(&#34;Incorrect file format.  Wrong number of columns.&#34;)
			}
			if _, ok := regions[fields[5]]; !ok {
				regions[fields[regionField]] = make([]customer, 0, 10)
			}
			regions[fields[regionField]] = append(regions[fields[regionField]], customer{num: fields[accountField], name: fields[nameField]})
			return nil
		})

	// don&#39;t continue if the wrong number of columns were observed
	if e != nil {
		panic(e)
	}

	return regions
}

// this function processes a &#39;job&#39;
// each job encapsulates the region to which it corresponds, the slice of customers in that region, and a channel for communicating potential matches
func process(j job) {
	// defer closing the job&#39;s channel
	defer close(j.comm)

	// do the LD comparisons and send those that are small enough (1/2 size) through the channel
	// the nested loop compares all pairs of accounts without comparing any pair twice
	var size int
	for i := 0; i &lt; len(j.accounts)-1; i&#43;&#43; {
		for k := i &#43; 1; k &lt; len(j.accounts); k&#43;&#43; {

			// compute the &#39;size&#39; of this  comparison, being the longer of the two strings
			il, kl := len(j.accounts[i].name), len(j.accounts[k].name)
			if il &gt; kl {
				size = il
			} else {
				size = kl
			}

			// compute the edit distance (Levenshtein)
			ed := ld(j.accounts[i].name, j.accounts[k].name)

			// if the edit distance is smaller than half the &#39;size&#39; computed above, send it through
			// this 1/2 ratio was chosen arbitrarily and is open to debate
			if ed &lt; (size / 2) {
				j.comm &lt;- []string{strconv.Itoa(ed), strconv.Itoa(size), j.accounts[i].num, j.accounts[i].name, j.accounts[k].num, j.accounts[k].name}
			}
		}
	}
}

func main() {
	// set the max number of processors Go can use equal to th enumber of logical CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// get the customer records from the input, organized by region
	regions := read(os.Stdin)

	// make a queue of jobs, where one job is one region in which to make comparisons
	q := ring.New(len(regions))

	// compute the number of comparisons to be performed and add a job for each region to the queue
	var comparisons int64
	for k, v := range regions {
		// this formula was derived from &#39;n choose 2&#39;
		// we are computing it for each region and summing them
		comparisons &#43;= int64(((len(v) * len(v)) - len(v)) / 2)

		// add a job to the queue and advance to the next Value
		q.Value = job{region: k, accounts: v, comm: make(chan []string, 10)}
		q = q.Next()
	}
	// using Stderr because Stdout will be dumped to the output file like so: cat &#34;filename.csv&#34; | crmld &gt; output.csv
	fmt.Fprintln(os.Stderr, &#34;Total number of comparisons to compute: &#34;, comparisons)

	// start a pool of active jobs
	p := q.Unlink(runtime.NumCPU())

	// start processing on the active jobs
	p.Do(
		func(j interface{}) {
			go process(j.(job))
		})

	// create a csv writer and an output buffer
	output := make([][]string, 0, bufSize)
	cw := csv.NewWriter(os.Stdout)

	// some counters for progress reporting
	var count int64

	// start writing the output
	for {
		// if there are no more jobs, we&#39;re done
		if p.Len() == 0 &amp;&amp; q.Len() == 0 {
			break
		}

		// try to read from the currently selected jobs&#39; channel
		// failing that, move on to check the next job
		select {
		case line, ok := &lt;-p.Value.(job).comm:
			if line != nil {
				output = append(output, line)
				if len(output) == cap(output) {
					for _, line := range output {
						if e := cw.Write(line); e != nil {
							panic(e)
						}
					}
					output = make([][]string, 0, bufSize)
				}
			} else if !ok {
				count &#43;= len(p.Value.(job).accounts)
				fmt.Fprintln(os.Stderr, &#34;Region %v completed: %v of %v in total&#34;, p.Value.(job).region, count, comparisons)
				p = p.Prev()
				p.Unlink(1)
				if q.Len() &gt; 0 {
					p.Link(q.Unlink(1))
					go process(p.Next().Value.(job))
				}
			}
			p = p.Next()
		default:
			p = p.Next()
		}
	}

	cw.Flush()
}
