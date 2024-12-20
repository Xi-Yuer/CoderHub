package imagerelationservicelogic

import (
	imageservicelogic "coderhub/rpc/coderhub/internal/logic/imageservice"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetImagesByEntityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetImagesByEntityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImagesByEntityLogic {
	return &GetImagesByEntityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetImagesByEntity 获取实体关联的图片列表
func (l *GetImagesByEntityLogic) GetImagesByEntity(in *coderhub.GetImagesByEntityRequest) (*coderhub.GetImagesByEntityResponse, error) {
	// 获取实体关联的图片列表
	batchGetImageService := imageservicelogic.NewBatchGetLogic(l.ctx, l.svcCtx)
	images, err := batchGetImageService.BatchGet(&coderhub.BatchGetRequest{
		ImageIds: []int64{in.EntityId},
	})
	if err != nil {
		return nil, err
	}
	// 将图片信息转换为ImageInfo
	imageInfos := make([]*coderhub.ImageInfo, 0, len(images.Images))
	for _, image := range images.Images {
		imageInfos = append(imageInfos, &coderhub.ImageInfo{
			ImageId:      image.ImageId,
			BucketName:   image.BucketName,
			ObjectName:   image.ObjectName,
			Url:          image.Url,
			ThumbnailUrl: image.ThumbnailUrl,
			ContentType:  image.ContentType,
			Size:         image.Size,
			Width:        image.Width,
			Height:       image.Height,
			UploadIp:     image.UploadIp,
			UserId:       image.UserId,
			CreatedAt:    image.CreatedAt,
		})
	}

	return &coderhub.GetImagesByEntityResponse{
		Images: imageInfos,
	}, nil
}
