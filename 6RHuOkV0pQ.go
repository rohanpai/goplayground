package main

import (
	&#34;fmt&#34;
	&#34;testing&#34;

	. &#34;github.com/juju/affinity&#34;
	&#34;github.com/juju/affinity/storage/mem&#34;
)

// MessageBoard defines the interface for a message board service.
type MessageBoard interface {
	// Lurk fetches a page of posts, by page number, lower = newer.
	Lurk(pageNum int) (string, error)
	// Read fetches a single post content.
	Read(id int) (string, error)
	// Post writes a message to the board, returns its id.
	Post(msg string) (int, error)
	// Sticky prevents the message from rolling off, keeps it at the top.
	Sticky(threadId int) error
	// Ban a user for some time.
	Ban(user User, seconds int) error
	// Delete a message.
	Delete(threadId int) error
}

var MessageBoardRoles RoleMap = NewRoleMap(LurkerRole, PosterRole, ModeratorRole)

var AccessDenied error = fmt.Errorf(&#34;Access denied&#34;)

// MessageBoardResource represents the entire message board.
var MessageBoardResource Resource = NewResource(&#34;message-board:&#34;,
	ReadPerm, ListPerm, PostPerm, DeletePerm, StickyPerm, BanPerm)

// mbConn is a connection to the message board service as a certain user.
type mbConn struct {
	*Access
	AsUser User
}

func (mb *mbConn) Lurk(pageNumber int) (string, error) {
	// Check that the user has list permissions on the message board
	can, err := mb.Can(mb.AsUser, ListPerm, MessageBoardResource)
	if err != nil {
		return &#34;&#34;, err
	}
	if !can {
		return &#34;&#34;, AccessDenied
	}
	// Get the page content
	return mb.loadPage(pageNumber)
}

func (mb *mbConn) Post(msg string) (int, error) {
	can, err := mb.Can(mb.AsUser, PostPerm, MessageBoardResource)
	if err != nil {
		return 0, err
	}
	if !can {
		return 0, AccessDenied
	}
	return mb.post(msg)
}

func main(t *testing.T) {
	// Let&#39;s set up an RBAC store. We&#39;ll use the in-memory store
	// for this example. You should use something more permanent like the Mongo store.
	store := mem.NewStore()
	// Admin lets us grant and revoke roles
	admin := NewAdmin(store, MessageBoardRoles)
	// Anonymous scheme users can lurk and that&#39;s all
	admin.Grant(User{Identity{&#34;anon&#34;, &#34;*&#34;}}, LurkerRole, MessageBoardResource)
	// Verified Gooble users can post
	admin.Grant(User{Identity{&#34;gooble&#34;, &#34;*&#34;}}, PosterRole, MessageBoardResource)

	// A wild anon appears
	anon := User{Identity{&#34;anon&#34;, &#34;10.55.61.128&#34;}}

	// Connect to the message board service as this user
	// In a web application, you&#39;ll likely derive the user from http.Request, using
	// OAuth, OpenID, cookies, etc.
	mb := &amp;mbConn{&amp;Access{store, MessageBoardRoles}, anon}

	// Print the first page of the message board. The MessageBoard will check
	// Access.Can(user, ListPerm, MessageBoardResource).
	content, err := mb.Lurk(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

	// A tame authenticated user appears. Reattach as tame user now.
	// In real life, this would likely be in a distinct http.Handler with its own session.
	tame := User{Identity{&#34;gooble&#34;, &#34;YourRealName&#34;}}
	mb = &amp;mbConn{&amp;Access{store, MessageBoardRoles}, tame}

	// Post a message.
	_, err = mb.Post(&#34;check &#39;em&#34;)
	if err != nil {
		panic(err)
	}
}
