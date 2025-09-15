package imageservicelogic

import (
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strconv"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Delete 删除图片
func (l *DeleteLogic) Delete(in *coderhub.DeleteRequest) (*coderhub.DeleteResponse, error) {
	// 权限校验
	var (
		userId string
		err    error
	)
	if userId, err = utils.GetUserMetaData(l.ctx); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	err = l.svcCtx.ImageRepository.Delete(l.ctx, in.ImageId)
	if err != nil {
		return nil, err
	}
	return &coderhub.DeleteResponse{
		Success: true,
	}, nil
}
