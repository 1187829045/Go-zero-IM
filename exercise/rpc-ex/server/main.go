package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type UserServer struct{}

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

func (u *UserServer) GetUser(req GetUsereq, rsp *GetUserrsp) error {
	if result, ok := Users[req.Id]; ok {
		*rsp = GetUserrsp{
			Id:    result.Id,
			Name:  result.Name,
			Phone: result.Phone,
		}
		return nil
	}
	return errors.New("用户不存在")
}
func main() {
	//创建好服务
	usersrv := new(UserServer)
	//服务注册到rpc里面
	rpc.Register(usersrv)
	//监听
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Println("启动成功")
	//连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)

	}
}
