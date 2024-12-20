package imageservicelogic

import (
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Get 获取图片信息
func (l *GetLogic) Get(in *coderhub.GetRequest) (*coderhub.ImageInfo, error) {
	imageModel, err := l.svcCtx.ImageRepository.GetByID(l.ctx, in.ImageId)
	if err != nil {
		return nil, err
	}
	return &coderhub.ImageInfo{
		ImageId:      imageModel.ID,
		BucketName:   imageModel.BucketName,
		ObjectName:   imageModel.ObjectName,
		Url:          imageModel.URL,
		ThumbnailUrl: imageModel.ThumbnailURL,
		ContentType:  imageModel.ContentType,
		Size:         imageModel.Size,
		Width:        imageModel.Width,
		Height:       imageModel.Height,
		UploadIp:     imageModel.UploadIP,
		UserId:       imageModel.UserID,
		CreatedAt:    imageModel.CreatedAt.Unix(),
	}, nil
}
