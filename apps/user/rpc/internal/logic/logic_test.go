/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package logic

import (
	"github.com/zeromicro/go-zero/core/conf"
	"llb-chat/apps/user/rpc/internal/config"
	"llb-chat/apps/user/rpc/internal/svc"
	"path/filepath"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
