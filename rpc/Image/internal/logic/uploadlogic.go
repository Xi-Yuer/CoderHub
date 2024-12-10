package logic

import (
	"bytes"
	"coderhub/model"
	"coderhub/rpc/Image/image"
	"coderhub/rpc/Image/internal/svc"
	"coderhub/shared/utils"
	"context"
	"fmt"
	"time"

	"math/rand"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Upload 上传图片
func (l *UploadLogic) Upload(in *image.UploadRequest) (*image.ImageInfo, error) {
	// 上传图片到 MinIO
	// 获取文件信息，这里的文件是一个字节切片
	fileReader := bytes.NewReader(in.File)
	// 获取文件名,这里没有文件名，所以使用时间戳 + 随机数生成一个文件名
	fileName := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(10000))

	// 上传图片并生成缩略图
	imageInfo, err := l.svcCtx.Minio.UploadImageWithThumbnail(l.svcCtx.Config.Minio.Bucket, fileName, fileReader, int64(len(in.File)), in.ContentType, l.svcCtx.Config.Minio.ThumbnailWidth)
	if err != nil {
		return nil, err
	}

	// 插入数据库
	ID := utils.GenID()
	err = l.svcCtx.ImageRepository.Create(l.ctx, &model.Image{
		ID:           ID,
		BucketName:   l.svcCtx.Config.Minio.Bucket,
		ObjectName:   fileName,
		URL:          imageInfo.URL,
		ThumbnailURL: imageInfo.ThumbnailURL,
		ContentType:  in.ContentType,
		Size:         imageInfo.Size,
		Width:        int32(imageInfo.Width),
		Height:       int32(imageInfo.Height),
		UploadIP:     "",
		UserID:       in.UserId,
		Status:       model.ImageStatusActive,
	})
	if err != nil {
		return nil, err
	}

	return &image.ImageInfo{
		ImageId:      ID,
		BucketName:   l.svcCtx.Config.Minio.Bucket,
		ObjectName:   fileName,
		Url:          imageInfo.URL,
		ThumbnailUrl: imageInfo.ThumbnailURL,
		ContentType:  in.ContentType,
		Size:         imageInfo.Size,
		Width:        int32(imageInfo.Width),
		Height:       int32(imageInfo.Height),
		UploadIp:     "",
		UserId:       in.UserId,
		Status:       model.ImageStatusActive,
		CreatedAt:    time.Now().Format(time.DateTime),
	}, nil
}
