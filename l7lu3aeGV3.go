package main

import (
	&#34;fmt&#34;
	&#34;strconv&#34;
	&#34;net/http&#34;
	&#34;github.com/julienschmidt/httprouter&#34;
	&#34;golang.org/x/net/context&#34;
)

// Unexported key type to prevent collision with context in other packages
type key int

// Key for the user ID
var UserIDKey key = 1

// Define a user
type User struct {
	ID    int
	Email string
}

type CtxHandle func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params)

func HandleUser(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := ctx.Value(UserIDKey).(*User)

	fmt.Fprintln(w, &#34;User ID: &#34; &#43; strconv.Itoa(user.ID))
	fmt.Fprintln(w, &#34;Email: &#34; &#43; user.Email)

	// Update the user
	user.ID = 5

	// Update context
	ctx = context.WithValue(ctx, UserIDKey, user)
}

func Middleware(next CtxHandle) CtxHandle {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := &amp;User{ID: 3, Email: &#34;test@email.com&#34;}

		ctx = context.WithValue(ctx, UserIDKey, user)
		next(ctx, w, r, ps)

		// Print updated user value
		user = ctx.Value(UserIDKey).(*User)
		fmt.Fprintln(w, &#34;User ID: &#34; &#43; strconv.Itoa(user.ID))
	}
}

func Wrapper(lead CtxHandle) httprouter.Handle {
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		lead(ctx, w, r, ps)
	}
}

func main() {
	router := httprouter.New()

	router.GET(&#34;/user&#34;, Wrapper(Middleware(HandleUser)))

	// Start the server
	fmt.Println(&#34;Listening...&#34;)
	http.ListenAndServe(&#34;:80&#34;, router)
}