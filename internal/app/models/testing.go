package models

import "testing"

//TestUser ...
func TestUser(t *testing.T) *User {
	return &User{
		Email:        "user@gmail.com",
		UserPassword: "password",
	}
}
