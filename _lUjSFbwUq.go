package user

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"dasa.cc/dae/context"
	"dasa.cc/dae/datastore"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id       bson.ObjectId `_id`
	Name     string
	Email    string
	Salt     []byte
	Password []byte
	DataBin  []byte
}

func (u *User) SetData(data interface{}) error {
	out, err := bson.Marshal(data)
	if err == nil {
		u.DataBin = out
	}
	return err
}

func (u *User) Data(out interface{}) error {
	return bson.Unmarshal(u.DataBin, out)
}

func New() *User {
	u := new(User)
	u.Id = bson.NewObjectId()
	return u
}

// GetSalt returns a new 512bit salt for securing a password.
func GetSalt() []byte {
	b := make([]byte, 64)
	// TODO inspect bytes read!
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return b
}

// Encrypt uses sha-512 to hash a given password and salt.
func Encrypt(password, salt []byte) []byte {
	h := sha512.New()
	h.Write(append(password, salt...))
	return h.Sum([]byte{})
}

// Validate checks the given password against the user's current password.
func (user *User) Validate(password string) bool {
	p := Encrypt([]byte(password), user.Salt)
	return bytes.Equal(p, user.Password)
}

// SetPassword updates the user's password on the struct, but does not save the changes automatically.
func (user *User) SetPassword(newPass string) {
	salt := GetSalt()
	pass := Encrypt([]byte(newPass), salt)

	user.Salt = salt
	user.Password = pass
}

func Current(c *context.Context, db *datastore.DB) *User {
	email := c.Session().Values["email"].(string)
	u, err := FindEmail(db, email)
	if err != nil {
		panic(err)
	}
	return u
}

func SetCurrent(c *context.Context, u *User) {
	c.Session().Values["auth"] = true
	c.Session().Values["email"] = u.Email
}

func DelCurrent(c *context.Context) {
	delete(c.Session().Values, "auth")
	delete(c.Session().Values, "email")
}

func FindEmail(db *datastore.DB, email string) (u *User, err error) {
	q := bson.M{"email": email}
	if err := db.C("users").Find(q).One(&u); err != nil {
		return nil, err
	}
	return u, nil
}
