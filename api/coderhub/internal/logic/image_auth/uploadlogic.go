package image_auth

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

// 上传图片
func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadLogic) Upload() (resp *types.UploadResponse, err error) {
	err = l.r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}
	file, handler, err := l.r.FormFile("file")
	if err != nil {
		return l.errorResp(err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			l.Logger.Error("关闭文件失败", err)
		}
	}(file)

	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return l.errorResp(err)
	}

	response, err := l.svcCtx.ImageAuthService.Upload(l.ctx, &coderhub.UploadRequest{
		File:        fileBytes,
		Filename:    handler.Filename,
		UserId:      userId,
		ContentType: handler.Header.Get("Content-Type"),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(response)
}
func (l *UploadLogic) successResp(response *coderhub.ImageInfo) (*types.UploadResponse, error) {
	return &types.UploadResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.ImageInfo{
			ImageId:      response.ImageId,
			Url:          response.Url,
			ThumbnailUrl: response.ThumbnailUrl,
			Width:        response.Width,
			Height:       response.Height,
		},
	}, nil
}
func (l *UploadLogic) errorResp(err error) (*types.UploadResponse, error) {
	return &types.UploadResponse{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
