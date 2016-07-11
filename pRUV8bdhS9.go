// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program demonstrating struct composition.
package main

import (
	&#34;errors&#34;
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;time&#34;
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// EOD represents the end of the data stream.
var EOD = errors.New(&#34;EOD&#34;)

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct{}

// Pull knows how to pull data out of Xenia.
func (Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return EOD

	case 5:
		return errors.New(&#34;Error reading data from Xenia&#34;)

	default:
		d.Line = &#34;Data&#34;
		fmt.Println(&#34;In:&#34;, d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct{}

// Store knows how to store data into Pillar.
func (Pillar) Store(d Data) error {
	fmt.Println(&#34;Out:&#34;, d.Line)
	return nil
}

// =============================================================================

// System wraps Xenia and Pillar together into a single system.
type System struct {
	Xenia
	Pillar
}

// =============================================================================

// IO provides support to copy bulk data.
type IO struct{}

// pull knows how to pull bulks of data from Xenia.
func (IO) pull(x *Xenia, data []Data) (int, error) {
	for i := range data {
		if err := x.Pull(&amp;data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data from Pillar.
func (IO) store(p *Pillar, data []Data) error {
	for _, d := range data {
		if err := p.Store(d); err != nil {
			return err
		}
	}

	return nil
}

// Copy knows how to pull and store data from the System.
func (io IO) Copy(sys *System, batch int) error {
	for {
		data := make([]Data, batch)

		i, err := io.pull(&amp;sys.Xenia, data)
		if i &gt; 0 {
			if err := io.store(&amp;sys.Pillar, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {

	// Initialize the system for use.
	sys := System{
		Xenia:  Xenia{},
		Pillar: Pillar{},
	}

	var io IO
	if err := io.Copy(&amp;sys, 3); err != EOD {
		fmt.Println(err)
	}
}
