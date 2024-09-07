package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/wuid"
	"llb-chat/pkg/xerr"
	"time"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 方法处理群组创建的逻辑
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: 在这里添加逻辑并删除此行

	// 创建新的群组对象
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource), // 生成群组唯一ID
		Name:       in.Name,                                       // 设置群组名称
		Icon:       in.Icon,                                       // 设置群组图标
		CreatorUid: in.CreatorUid,                                 // 设置群组创建者ID
		//IsVerify:   true,         // 是否需要验证，暂时注释
		IsVerify: false, // 设置群组不需要验证
	}

	// 开始事务
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务中插入群组记录
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)
		if err != nil {
			// 如果插入群组记录时出错，返回包装后的数据库错误
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		// 插入群组成员记录，创建者为群主
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,                            // 设置群组ID
			UserId:    in.CreatorUid,                        // 设置用户ID为创建者ID
			RoleLevel: int(constants.CreatorGroupRoleLevel), // 设置角色级别为群主
		})
		if err != nil {
			// 如果插入群组成员记录时出错，返回包装后的数据库错误
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil // 事务成功结束
	})

	// 暂停2秒，模拟处理时间
	time.Sleep(2 * time.Second)

	// 成功创建群组，返回群组ID和错误信息
	return &social.GroupCreateResp{
		Id: groups.Id, // 返回创建的群组ID
	}, err // 返回错误信息（如果有）
}
