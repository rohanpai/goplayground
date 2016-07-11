// Sample program to show the basic concept of using a pointer
// to share data.
package main

import &#34;fmt&#34;

// user represents a user in the system.
type user struct {
	name   string
	email  string
	logins int
}

// main is the entry point for the application.
func main() {
	// Declare and initialize a variable named bill of type user.
	bill := user{
		name:  &#34;Bill&#34;,
		email: &#34;bill@ardanstudios.com&#34;,
	}

	//** We don&#39;t need to include all the fields when specifying field
	// names with a struct literal.

	// Pass the &#34;address of&#34; the bill value.
	display(&amp;bill)

	// Pass the &#34;address of&#34; the logins field from within the bill value.
	increment(&amp;bill.logins)

	// Pass the &#34;address of&#34; the bill value.
	display(&amp;bill)
}

// increment declares logins as a pointer variable whose value is
// always an address and points to values of type int.
func increment(logins *int) {
	*logins&#43;&#43;
	fmt.Printf(&#34;&amp;logins[%p] logins[%p] *logins[%d]\n&#34;, &amp;logins, logins, *logins)
}

// display declares u as user pointer variable whose value is always an address
// and points to values of type user.
func display(u *user) {
	fmt.Printf(&#34;%p\t%&#43;v\n&#34;, u, *u)
}
