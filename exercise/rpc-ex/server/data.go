package main

type User struct {
	Id    string
	Name  string
	Phone string
}

var Users = map[string]*User{
	"1": {
		Id:    "1",
		Name:  "木兮",
		Phone: "123456",
	},
	"2": {
		Id:    "2",
		Name:  "木兮2",
		Phone: "654321",
	},
}
