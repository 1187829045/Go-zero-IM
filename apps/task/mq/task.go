/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"llb-chat/apps/task/mq/internal/config"
	"llb-chat/apps/task/mq/internal/handler"
	"llb-chat/apps/task/mq/internal/svc"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	serviceGroup.Start()
}
