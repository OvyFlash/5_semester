package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

//User ...
type User struct {
	UserID            int32  `json:"userid,omitempty"`
	UserName          string `json:"username,omitempty"`
	FirstName         string `json:"firstname"`
	LastName          string `json:"lastname,omitempty"`
	Email             string `json:"email"`
	UserPassword      string `json:"userpassword,omitempty"`
	PhoneNumber       int    `json:"phonenumber,omitempty"`
	EncryptedPassword string `json:"-"`
	JSONToken         string `json:"jwt,omitempty"`
	ProfilePic        string `json:"pic_url,omitempty"`
}

//UserDecode user for decoding
type UserDecode struct {
	UserID            int32          `json:"userid"`
	UserName          sql.NullString `json:"username,omitempty"`
	FirstName         string         `json:"firstname"`
	LastName          sql.NullString `json:"lastname,omitempty"`
	Email             string         `json:"email"`
	PhoneNumber       sql.NullInt32  `json:"phonenumber,omitempty"`
	EncryptedPassword string         `json:"-"`
}

//ToUser ...
func (u *UserDecode) ToUser() *User {
	user := &User{
		UserID:    u.UserID,
		Email:     u.Email,
		FirstName: u.FirstName,
	}
	if u.UserName.Valid {
		user.UserName = u.UserName.String
	}
	if u.LastName.Valid {
		user.LastName = u.LastName.String
	}
	if u.PhoneNumber.Valid {
		user.PhoneNumber = int(u.PhoneNumber.Int32)
	}

	return user
}
func (u *User) String() string {
	usr := "USER\n"
	if u.UserID != 0 {
		usr += fmt.Sprintf("UserID: %d\n", u.UserID)
	}
	if u.UserName != "" {
		usr += fmt.Sprintf("UserName: %s\n", u.UserName)
	}
	if u.FirstName != "" {
		usr += fmt.Sprintf("FirstName: %s\n", u.FirstName)
	}
	if u.LastName != "" {
		usr += fmt.Sprintf("LastName: %s\n", u.LastName)
	}
	if u.Email != "" {
		usr += fmt.Sprintf("Email: %s\n", u.Email)
	}
	// if u.UserPassword != "" {
	// 	usr += "Password: +\n"
	// } else {
	// 	usr += "Password: -\n"
	// }
	if u.PhoneNumber != 0 {
		usr += fmt.Sprintf("PhoneNumber: %d\n", u.PhoneNumber)
	}
	usr += "\n"
	return usr
}

//CheckEmailAndPassword ...
func (u *User) CheckEmailAndPassword() error {
	if u.Email == "" || u.FirstName == "" {
		return errors.New("Forgot email or firstname?")
	}
	if u.UserPassword == "" && u.EncryptedPassword == "" {
		return errors.New("Forgot password?")
	}
	if len(u.UserPassword) < 6 && u.EncryptedPassword == "" {
		return errors.New("Password must contain 8 or more symbols")
	}
	re := regexp.MustCompile(`(?:[a-z0-9!#$%&'*+\/=?^_\x60{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_\x60{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)
	if string(re.Find([]byte(u.Email))) == "" {
		return errors.New("Bad email")
	}
	return nil
}

//EncryptPassword ...
func (u *User) EncryptPassword() error {
	b, err := bcrypt.GenerateFromPassword([]byte(u.UserPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(b)
	return nil
}

//ClearPassword ...
func (u *User) ClearPassword() {
	u.UserPassword = ""
}

//ComparePassword compares encrypted password and password that user has entered
func (u *User) ComparePassword(user *User) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(user.UserPassword)) == nil
}

//Difference compares new user with old in DB and changes values
func (u *User) Difference(user *User) {
	if u.UserName != user.UserName && user.UserName != "" {
		u.UserName = user.UserName
	}
	if u.FirstName != user.FirstName && user.FirstName != "" {
		u.FirstName = user.FirstName
	}
	if u.LastName != user.LastName && user.LastName != "" {
		u.LastName = user.LastName
	}
	if u.Email != user.Email && user.Email != "" {
		u.Email = user.Email
	}
	if !u.ComparePassword(user) && user.UserPassword != "" {
		u.UserPassword = user.UserPassword
	}
	if u.PhoneNumber != user.PhoneNumber && user.PhoneNumber != 0 {
		u.PhoneNumber = user.PhoneNumber
	}
}

/*
CREATE TABLE IF NOT EXISTS Users (
    userid BIGINT NOT NULL AUTO_INCREMENT,
    username text,
    firstname text,
    lastname text,
    email varchar(512) not null,
    encrypted_password text not null,
    phone_number int,
    UNIQUE(email),
    PRIMARY KEY(userid)
);*/

//UserDecode exists to avoid sql null errors
