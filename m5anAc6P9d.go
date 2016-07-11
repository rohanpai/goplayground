package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
type Person interface {
	GetReady(*sync.WaitGroup)
	PutOnShoes(*sync.WaitGroup)
}
*/

type person string

func (p person) GetReady(wg *sync.WaitGroup) {
	fmt.Println(string(p), "started to get ready.")
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(t)
		fmt.Println(string(p), "spent", t, "getting ready.")

	}()
}

func (p person) PutOnShoes(wg *sync.WaitGroup) {
	fmt.Println(string(p), "started putting on shoes.")
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(t)
		fmt.Println(string(p), "spent", t, "putting on shoes.")

	}()
}

type Alarm interface {
	ArmTheAlarm() <-chan struct{}
}

type alarm struct{}

func (alarm) ArmTheAlarm() <-chan struct{} {
	fmt.Println("Arming the alarm.")
	done := make(chan struct{})

	go func() {
		fmt.Println("Alarm is counting down.")
		time.Sleep(5 * time.Second)
		fmt.Println("Alarm is armed!")
		done <- struct{}{}
	}()
	return done
}

func main() {
	rand.Seed(0)

	wg := new(sync.WaitGroup)
	alice, bob := person("Alice"), person("Bob")

	fmt.Println("Let's go for a walk!")

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

	fmt.Println("Exiting and locking the door.")

	// Wait for the alarm to arm itself
	<-done
}
