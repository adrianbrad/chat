package users

import "github.com/adrianbrad/chat/auth"

var Users = map[string]*auth.User{
	"1": &auth.User{Name: "brad", Role: true},
	"2": &auth.User{Name: "john", Role: false},
	"3": &auth.User{Name: "eusebiu", Role: true},
}
