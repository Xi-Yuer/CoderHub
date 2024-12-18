package logic

import (
	"context"
	"time"

	"coderhub/rpc/Image/image"
	"coderhub/rpc/Image/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListByUserLogic {
	return &ListByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListByUser 获取用户图片列表
func (l *ListByUserLogic) ListByUser(in *image.ListByUserRequest) (*image.ListByUserResponse, error) {
	images, total, err := l.svcCtx.ImageRepository.ListByUserID(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	imageInfos := make([]*image.ImageInfo, 0, len(images))
	for _, v := range images {
		imageInfos = append(imageInfos, &image.ImageInfo{
			ImageId:      v.ID,
			BucketName:   v.BucketName,
			ObjectName:   v.ObjectName,
			Url:          v.URL,
			ThumbnailUrl: v.ThumbnailURL,
			ContentType:  v.ContentType,
			Size:         v.Size,
			Width:        v.Width,
			Height:       v.Height,
			UploadIp:     v.UploadIP,
			UserId:       v.UserID,
			Status:       v.Status,
			CreatedAt:    v.CreatedAt.Format(time.DateTime),
		})
	}

	return &image.ListByUserResponse{
		Images:  imageInfos,
		Total:   total,
	}, nil
}
