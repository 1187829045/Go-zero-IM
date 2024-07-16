package main

import (
	"log"
	"net/rpc"
)

type (
	GetUsereq struct {
		Id string
	}
	GetUserrsp struct {
		Id    string
		Name  string
		Phone string
	}
)

func main() {
	//建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("请求失败")
	}
	defer client.Close()
	//调用请求
	var (
		req = GetUsereq{Id: "1"}
		rsp GetUserrsp
	)
	err = client.Call("UserServer.GetUser", req, &rsp)
	if err != nil {
		log.Println("调用出错")
	}
	log.Println(rsp)
}
