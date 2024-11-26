// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: image.proto

package server

import (
	"context"

	"coderhub/rpc/Image/image"
	"coderhub/rpc/Image/internal/logic"
	"coderhub/rpc/Image/internal/svc"
)

type ImageServer struct {
	svcCtx *svc.ServiceContext
	image.UnimplementedImageServer
}

func NewImageServer(svcCtx *svc.ServiceContext) *ImageServer {
	return &ImageServer{
		svcCtx: svcCtx,
	}
}

// 上传图片
func (s *ImageServer) Upload(ctx context.Context, in *image.UploadRequest) (*image.ImageInfo, error) {
	l := logic.NewUploadLogic(ctx, s.svcCtx)
	return l.Upload(in)
}

// 删除图片
func (s *ImageServer) Delete(ctx context.Context, in *image.DeleteRequest) (*image.DeleteResponse, error) {
	l := logic.NewDeleteLogic(ctx, s.svcCtx)
	return l.Delete(in)
}

// 获取图片信息
func (s *ImageServer) Get(ctx context.Context, in *image.GetRequest) (*image.ImageInfo, error) {
	l := logic.NewGetLogic(ctx, s.svcCtx)
	return l.Get(in)
}

// 获取用户图片列表
func (s *ImageServer) ListByUser(ctx context.Context, in *image.ListByUserRequest) (*image.ListByUserResponse, error) {
	l := logic.NewListByUserLogic(ctx, s.svcCtx)
	return l.ListByUser(in)
}
