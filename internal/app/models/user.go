package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

//User ...
type User struct {
	UserID            int32  `json:"userid"`
	UserName          string `json:"username,omitempty"`
	FirstName         string `json:"firstname,omitempty"`
	LastName          string `json:"lastname,omitempty"`
	Email             string `json:"email"`
	UserPassword      string `json:"userpassword,omitempty"`
	PhoneNumber       int    `json:"phonenumber,omitempty"`
	EncryptedPassword string `json:"-"`
}

//Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Email, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)))
	//user should be valid when returned from DATABASE
	//When we get user from database we dont know about his simple password (u.UserPassword == "")
	//If encrypted password is not nil then user is valid
	//If UserPassword is not nil then user is valid

}

//BeforeCreate all we need to do before user creating
func (u *User) BeforeCreate() error {
	//firstName, email, password
	if len(u.UserPassword) > 0 {
		enc, err := encryptString(u.UserPassword)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

//Sanitize delete password
func (u *User) Sanitize() {
	u.UserPassword = ""
}

//ComparePassword compares encrypted password and password that user has entered
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
