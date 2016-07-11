package main

import (
	"bson"
	"labix.org/v2/mgo"
)

type User struct {
	Email   string `bson:"email"`
	Name    string `bson:"name"`
	Enabled bool   `bson:"enabled"`
	GroupId int    `bson:group_id"`
}

func update_user(w http.RequestWriter, r *http.Request) {
	user_id := "123"

	//get the updates to the user object
	json_data := `{ "email":"user@example.com", "name":"My New Name"}`
	new_user := &User{}
	_ = json.Unmarshal(json_data, new_user)

	//connect to MongoDB
	session, _ := mgo.Dial("mongodb://localhost/test")
	col := session.DB("test").C("users")

	//search for an existing User by e-mail
	existing_user := &User{}
	_ = col.FindId(user_id).One(existing_user)

	//...
	//assuming the user was found:

	_ = col.UpdateId(user_id, new_user)

	/* oops, new group id is now 0, enabled is 0 */
}
