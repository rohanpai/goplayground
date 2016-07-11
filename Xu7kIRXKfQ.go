package main

import (
	&#34;bson&#34;
	&#34;labix.org/v2/mgo&#34;
)

type User struct {
	Email   string `bson:&#34;email&#34;`
	Name    string `bson:&#34;name&#34;`
	Enabled bool   `bson:&#34;enabled&#34;`
	GroupId int    `bson:group_id&#34;`
}

func update_user(w http.RequestWriter, r *http.Request) {
	user_id := &#34;123&#34;

	//get the updates to the user object
	json_data := `{ &#34;email&#34;:&#34;user@example.com&#34;, &#34;name&#34;:&#34;My New Name&#34;}`
	new_user := &amp;User{}
	_ = json.Unmarshal(json_data, new_user)

	//connect to MongoDB
	session, _ := mgo.Dial(&#34;mongodb://localhost/test&#34;)
	col := session.DB(&#34;test&#34;).C(&#34;users&#34;)

	//search for an existing User by e-mail
	existing_user := &amp;User{}
	_ = col.FindId(user_id).One(existing_user)

	//...
	//assuming the user was found:

	_ = col.UpdateId(user_id, new_user)

	/* oops, new group id is now 0, enabled is 0 */
}
