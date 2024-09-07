package main

type User struct {
	Id    string
	Name  string
	Phone string
}

var users = map[string]*User{
	"1": {
		Id:    "1",
		Name:  "木兮",
		Phone: "13800001111",
	},
	"2": {
		Id:    "2",
		Name:  "小慕",
		Phone: "15688880000",
	},
}
