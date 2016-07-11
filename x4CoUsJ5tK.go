package main

// sort a map&#39;s keys in descending order of its values.

import &#34;sort&#34;
import &#34;fmt&#34;

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] &gt; sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i&#43;&#43;
	}
	sort.Sort(sm)
	return sm.s
}

func main() {
	s := []string{&#34;Python&#34;, &#34;Python&#34;, &#34;Python&#34;, &#34;igor&#34;, &#34;igor&#34;, &#34;igor&#34;, &#34;igor&#34;, &#34;go&#34;, &#34;go&#34;, &#34;Golang&#34;, &#34;Golang&#34;, &#34;Golang&#34;, &#34;Golang&#34;, &#34;Py&#34;, &#34;Py&#34;}
	count := make(map[string]int)

	for _, v := range s {
		count[v]&#43;&#43;
	}
	
	for _, res := range sortedKeys(count) {
		fmt.Println(res, count[res])
	}

}

