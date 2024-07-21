/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"llb-chat/apps/im/ws/internal/config"
	"llb-chat/apps/im/ws/internal/handler"
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/apps/im/ws/websocket"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		//websocket.WithServerAck(websocket.RigorAck),
		//websocket.WithServerMaxConnectionIdle(10*time.Second),
	)
	defer srv.Stop()
	handler.RegisterHandlers(srv, ctx)
	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")
	srv.Start()
}
