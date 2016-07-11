package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// doing hw 4.6 of CS215
// made good headway learned the csv package, bufio
// bunch of other file opening, manipulations etc
// gonna try to make it rly good with go's concurrency

// This is represent the bipartite graph
type Graph map[string]map[string]int

type Imdb struct {
	actor string
	score float64
}

type ImdbRanks []Imdb

func (s ImdbRanks) Len() int           { return len(s) }
func (s ImdbRanks) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ImdbRanks) Less(i, j int) bool { return s[i].score < s[j].score }

func g1(ch chan<- int) {
	for i := 2; ; i = i + 2 {
		ch <- i // Send 'i' to channel 'ch'
		fmt.Printf("sending %d from g1\n", i)
	}
}

func g2(ch chan<- int) {
	for i := 2; ; i = i + 3 {
		ch <- i // Send 'i' to channel 'ch'
		fmt.Printf("sending %d from g2\n", i)
	}
}

func make_link(g Graph, actor, movie string) {
	if _, in := g[actor]; !in {
		g[actor] = make(map[string]int)
	}
	g[actor][movie] = 1

	if _, in := g[movie]; !in {
		g[movie] = make(map[string]int)
	}
	g[movie][actor] = 1
}

// this is causing problems!!! g is a bipartite graph
func average_centrality(g Graph, node string) (dis float64) {
	dis_from_start := map[string]int{node: 0}
	open_list := []string{node}
	for len(open_list) != 0 {
		current := open_list[0]
		open_list = open_list[1:]
		for neighbor := range g[current] {
			if _, ok := dis_from_start[neighbor]; !ok {
				dis_from_start[neighbor] = dis_from_start[current] + 1
				open_list = append(open_list, neighbor)
			}
		}
	}
	for _, v := range dis_from_start {
		dis += float64(v)
	}
	dis = dis / float64(len(dis_from_start))
	return
}

func test() float64 {

	time.Sleep(100 * time.Millisecond)
	return 0.0
}

func main() {
	t1 := time.Now()
	fp, err := os.Open("/Users/dluna/Downloads/file.tsv")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer fp.Close()
	r := csv.NewReader(fp)
	r.Comma, r.FieldsPerRecord = '\t', 3

	all_recs, err := r.ReadAll() // gives a 2d slice of the file r, each slice is a line split into 3 sections
	if err != nil {
		fmt.Println(err.Error())
	}
	g := make(Graph)
	actors := make(map[string]bool) // create a map of actors, used later on to compute centralities
	for _, line := range all_recs {
		make_link(g, line[0], line[1]+" "+line[2])
		actors[line[0]] = true
	}

	// make the slice, ch and wg
	top20k := make(ImdbRanks, len(actors))
	wg := &sync.WaitGroup{}
	i := 0
	// find centrality concurrently
	for node, _ := range actors {
		wg.Add(1)
		go func(arr ImdbRanks, i int, node string) {
			defer wg.Done()
			sc := average_centrality(g, node)
			// sc := test() // func to test is sleep blocks, this is not the case
			arr[i] = Imdb{node, sc}
		}(top20k, i, node)

		i++
		// when running average_centrality it shows 3 goroutines up at the same time
		// when running test() number goes up all the way to len(actors)
		// fmt.Println(wg, i, runtime.NumGoroutine())
	}

	wg.Wait()

	// unload all the Imdb into the array
	// for i := 0; i < len(actors); i++ {
	// 	top20k[i] = <-ch
	// 	fmt.Println(top20k[i])
	// }

	sort.Sort(top20k)

	for i := 0; i < 20; i++ {
		fmt.Println(top20k[i].actor, top20k[i].score)
	}
	fmt.Println(time.Since(t1))
}
