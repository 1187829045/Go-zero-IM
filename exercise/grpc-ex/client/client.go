package main

import (
	"Go-zero-IM/exercise/3-6/proto/user"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	client, err := grpc.Dial("127.0.0.1:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("连接失败", err)
	}
	defer client.Close()

	c := user.NewUserClient(client)

	r, err := c.GetUser(context.Background(), &user.GetUserReq{Id: "1"})
	if err != nil {
		log.Fatal("请求失败", err)
	}

	log.Println(r)
}
