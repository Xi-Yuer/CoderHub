package questions_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateQuestionBankLogic 创建题库
func NewCreateQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionBankLogic {
	return &CreateQuestionBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionBankLogic) CreateQuestionBank(req *types.CreateQuestionBankReq) (resp *types.CreateQuestionBankResp, err error) {
	UserId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	imageInfo, err := l.svcCtx.ImagesService.Get(l.ctx, &coderhub.GetRequest{
		ImageId: req.CoverImage,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.QuestionBankService.CreateQuestionBank(l.ctx, &coderhub.CreateQuestionBankRequest{
		Name:        req.Name,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
		CreateUser:  UserId,
		CoverImage: &coderhub.Image{
			ImageId:      imageInfo.ImageId,
			Url:          imageInfo.Url,
			ThumbnailUrl: imageInfo.ThumbnailUrl,
			Width:        imageInfo.Width,
			Height:       imageInfo.Height,
		},
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *CreateQuestionBankLogic) errorResp(err error) (resp *types.CreateQuestionBankResp, err1 error) {
	return &types.CreateQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateQuestionBankLogic) successResp() (resp *types.CreateQuestionBankResp, err1 error) {
	return &types.CreateQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "创建题库成功",
		},
		Data: true,
	}, nil
}
