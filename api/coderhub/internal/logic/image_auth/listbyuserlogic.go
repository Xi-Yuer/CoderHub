package image_auth

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListByUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户图片列表
func NewListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListByUserLogic {
	return &ListByUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListByUserLogic) ListByUser(req *types.ListByUserRequest) (resp *types.ListByUserResponse, err error) {
	response, err := l.svcCtx.ImageAuthService.ListByUser(l.ctx, &coderhub.ListByUserRequest{
		UserId:   req.UserId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(response)
}
func (l *ListByUserLogic) successResp(response *coderhub.ListByUserResponse) (*types.ListByUserResponse, error) {
	list := types.ImageInfoList{
		List:  make([]types.ImageInfo, 0, len(response.Images)),
		Total: response.Total,
	}
	for _, val := range response.Images {
		list.List = append(list.List, types.ImageInfo{
			ImageId:      val.ImageId,
			BucketName:   val.BucketName,
			ObjectName:   val.ObjectName,
			Url:          val.Url,
			ThumbnailUrl: val.ThumbnailUrl,
			ContentType:  val.ContentType,
			Size:         val.Size,
			Width:        val.Width,
			Height:       val.Height,
			UploadIp:     val.UploadIp,
			UserId:       val.UserId,
			CreatedAt:    val.CreatedAt,
		})
	}
	return &types.ListByUserResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &list,
	}, nil
}
func (l *ListByUserLogic) errorResp(err error) (*types.ListByUserResponse, error) {
	return &types.ListByUserResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
