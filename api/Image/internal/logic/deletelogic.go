package logic

import (
	"coderhub/shared/utils"
	"context"

	"coderhub/api/Image/internal/svc"
	"coderhub/api/Image/internal/types"
	"coderhub/conf"
	"coderhub/rpc/Image/image"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteLogic 删除图片
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteRequest) (resp *types.DeleteResponse, err error) {
	// 权限校验
	_, err = utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	ctx := utils.SetUserMetaData(l.ctx) // 设置元数据

	_, err = l.svcCtx.ImageService.Delete(ctx, &image.DeleteRequest{
		ImageId: req.ImageId,
		UserId:  req.UserId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteLogic) successResp() (*types.DeleteResponse, error) {
	return &types.DeleteResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteLogic) errorResp(err error) (*types.DeleteResponse, error) {
	return &types.DeleteResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
