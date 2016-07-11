// Sample program to show how polymorphic behavior with interfaces.
package main

import &#34;fmt&#34;

// notifier is an interface that defines notification
// type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements the notifier interface with a pointer receiver.
func (u *user) notify() {
	fmt.Printf(&#34;Sending user Email To %s&lt;%s&gt;\n&#34;,
		u.name,
		u.email)
}

// admin defines a admin in the program.
type admin struct {
	name  string
	email string
}

// notify implements the notifier interface with a pointer receiver.
func (a *admin) notify() {
	fmt.Printf(&#34;Sending admin Email To %s&lt;%s&gt;\n&#34;,
		a.name,
		a.email)
}

// main is the entry point for the application.
func main() {
	// Create two values one of type user and one of type admin.
	bill := user{&#34;Bill&#34;, &#34;bill@email.com&#34;}
	jill := admin{&#34;Jill&#34;, &#34;jill@email.com&#34;}

	// Pass a pointer of the values to support the interface.
	sendNotification(&amp;bill)
	sendNotification(&amp;jill)
}

// sendNotification accepts values that implement the notifier
// interface and sends notifications.
func sendNotification(n notifier) {
	n.notify()
}
