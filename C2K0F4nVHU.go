package main

import (
	&#34;fmt&#34;
	&#34;regexp&#34;
)

func main() {
	str := []byte(&#34;12:00&#34;)
	assigned := regexp.MustCompile(&#34;(.*):(.*)&#34;)
	group := assigned.FindSubmatch(str)
	fmt.Println(string(group[0]));
	fmt.Println();
	fmt.Println(string(group[1]))
	fmt.Println(string(group[2]))
}
