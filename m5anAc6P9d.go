package main

import (
	&#34;fmt&#34;
	&#34;math/rand&#34;
	&#34;sync&#34;
	&#34;time&#34;
)

/*
type Person interface {
	GetReady(*sync.WaitGroup)
	PutOnShoes(*sync.WaitGroup)
}
*/

type person string

func (p person) GetReady(wg *sync.WaitGroup) {
	fmt.Println(string(p), &#34;started to get ready.&#34;)
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(t)
		fmt.Println(string(p), &#34;spent&#34;, t, &#34;getting ready.&#34;)

	}()
}

func (p person) PutOnShoes(wg *sync.WaitGroup) {
	fmt.Println(string(p), &#34;started putting on shoes.&#34;)
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(t)
		fmt.Println(string(p), &#34;spent&#34;, t, &#34;putting on shoes.&#34;)

	}()
}

type Alarm interface {
	ArmTheAlarm() &lt;-chan struct{}
}

type alarm struct{}

func (alarm) ArmTheAlarm() &lt;-chan struct{} {
	fmt.Println(&#34;Arming the alarm.&#34;)
	done := make(chan struct{})

	go func() {
		fmt.Println(&#34;Alarm is counting down.&#34;)
		time.Sleep(5 * time.Second)
		fmt.Println(&#34;Alarm is armed!&#34;)
		done &lt;- struct{}{}
	}()
	return done
}

func main() {
	rand.Seed(0)

	wg := new(sync.WaitGroup)
	alice, bob := person(&#34;Alice&#34;), person(&#34;Bob&#34;)

	fmt.Println(&#34;Let&#39;s go for a walk!&#34;)

	// Get ready
	bob.GetReady(wg)
	alice.GetReady(wg)
	wg.Wait()

	// Set the alarm
	done := new(alarm).ArmTheAlarm()

	// Put on shoes
	alice.PutOnShoes(wg)
	bob.PutOnShoes(wg)
	wg.Wait()

	fmt.Println(&#34;Exiting and locking the door.&#34;)

	// Wait for the alarm to arm itself
	&lt;-done
}
