package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
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

	fmt.Fprintln(w, "User ID: " + strconv.Itoa(user.ID))
	fmt.Fprintln(w, "Email: " + user.Email)

	// Update the user
	user.ID = 5

	// Update context
	ctx = context.WithValue(ctx, UserIDKey, user)
}

func Middleware(next CtxHandle) CtxHandle {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := &User{ID: 3, Email: "test@email.com"}

		ctx = context.WithValue(ctx, UserIDKey, user)
		next(ctx, w, r, ps)

		// Print updated user value
		user = ctx.Value(UserIDKey).(*User)
		fmt.Fprintln(w, "User ID: " + strconv.Itoa(user.ID))
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

	router.GET("/user", Wrapper(Middleware(HandleUser)))

	// Start the server
	fmt.Println("Listening...")
	http.ListenAndServe(":80", router)
}