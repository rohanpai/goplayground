package main

import (
	&#34;fmt&#34;
	&#34;time&#34;
)

func debounce(interval time.Duration, done chan&lt;- struct{}, input &lt;-chan int, f func(arg int)) {
	var (
		item int
		timeout = time.After(interval)
	)
	for {
		select {
		case item = &lt;-input:
			if item % 100000 == 0 {
				fmt.Println(&#34;received a send on a spammy channel - might be doing a costly operation if not for debounce&#34;)
			}
		case &lt;-timeout:
			f(item)
			done&lt;-struct{}{}
		}
	}
}

func main() {
	spammyChan := make(chan int, 10)
	done := make(chan struct{})
	go debounce(100*time.Millisecond, done, spammyChan, func(arg int) {
		fmt.Println(&#34;*****************************&#34;)
		fmt.Println(&#34;* DOING A COSTLY OPERATION! *&#34;)
		fmt.Println(&#34;*****************************&#34;)
		fmt.Println(&#34;In case you were wondering, the value passed to this function is&#34;, arg)
		fmt.Println(&#34;We could have more args to our \&#34;compiled\&#34; debounced function too, if we wanted.&#34;)
	})
	loop:
	for i := 0; i &lt; 10000000; i&#43;&#43; {
		select {
			case spammyChan &lt;- i:
				continue loop
			case &lt;-done:
				fmt.Println(&#34;break loop&#34;)
				return
		}
	}
	fmt.Println(&#34;exited loop&#34;)
	time.Sleep(500 * time.Millisecond)
	fmt.Println(&#34;done.&#34;)
}