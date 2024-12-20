package imageservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetLogic {
	return &BatchGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchGet 批量获取图片信息
func (l *BatchGetLogic) BatchGet(in *coderhub.BatchGetRequest) (*coderhub.BatchGetResponse, error) {
	images, err := l.svcCtx.ImageRepository.BatchGetImagesByID(l.ctx, in.ImageIds)
	if err != nil {
		return nil, err
	}
	imageInfos := make([]*coderhub.UploadImageInfo, 0)
	for _, value := range images {
		imageInfos = append(imageInfos, &coderhub.UploadImageInfo{
			ImageId:      value.ID,
			BucketName:   value.BucketName,
			ObjectName:   value.ObjectName,
			Url:          value.URL,
			ThumbnailUrl: value.ThumbnailURL,
			ContentType:  value.ContentType,
			Size:         value.Size,
			Width:        value.Width,
			Height:       value.Height,
			UploadIp:     value.UploadIP,
			UserId:       value.UserID,
			Status:       value.Status,
			CreatedAt:    value.CreatedAt.Unix(),
		})
	}

	return &coderhub.BatchGetResponse{
		Images: imageInfos,
	}, nil
}
