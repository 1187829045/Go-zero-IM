package main

import (
	"Go-zero-IM/exercise/3-6/proto/user"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

type UserServer struct{}

func (u *UserServer) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserResp, error) {
	if u1, ok := users[req.Id]; ok {
		return &user.GetUserResp{
			Id:    u1.Id,
			Name:  u1.Name,
			Phone: u1.Phone,
		}, nil
	}

	return nil, errors.New("不存在查询用户")
}

func main() {
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("监听失败", err)
	}
	s := grpc.NewServer()

	user.RegisterUserServer(s, new(UserServer))

	log.Println("服务已启动")

	s.Serve(listen)
}
