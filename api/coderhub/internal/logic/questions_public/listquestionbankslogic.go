package questions_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuestionBanksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListQuestionBanksLogic 获取题库列表
func NewListQuestionBanksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionBanksLogic {
	return &ListQuestionBanksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuestionBanksLogic) ListQuestionBanks(req *types.GetQuestionBankListReq) (resp *types.GetQuestionBankListResp, err error) {
	list, err := l.svcCtx.QuestionBankService.GetQuestionBankList(l.ctx, &coderhub.GetQuestionBankListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	banks := make([]*types.QuestionBank, 0, len(list.Banks))
	for _, bank := range list.Banks {
		banks = append(banks, &types.QuestionBank{
			Id:          bank.Id,
			Name:        bank.Name,
			Description: bank.Description,
			Difficulty:  bank.Difficulty,
			Tags:        bank.Tags,
			CoverImage: &types.ImageInfo{
				ImageId:      bank.CoverImage.ImageId,
				BucketName:   bank.CoverImage.BucketName,
				ObjectName:   bank.CoverImage.ObjectName,
				Url:          bank.CoverImage.Url,
				ThumbnailUrl: bank.CoverImage.ThumbnailUrl,
				ContentType:  bank.CoverImage.ContentType,
				Size:         bank.CoverImage.Size,
				Width:        bank.CoverImage.Width,
				Height:       bank.CoverImage.Height,
				UploadIp:     bank.CoverImage.UploadIp,
				UserId:       bank.CoverImage.UserId,
				CreatedAt:    bank.CoverImage.CreatedAt,
			},
			CreateUser: nil,
			CreatedAt:  bank.CreateTime,
			UpdatedAt:  bank.UpdateTime,
		})
	}
	return l.successResp(banks, list.Total)
}

func (l *ListQuestionBanksLogic) errorResp(err error) (*types.GetQuestionBankListResp, error) {
	return &types.GetQuestionBankListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.QuestionBankList{
			Total: 0,
			List:  nil,
		},
	}, nil
}
func (l *ListQuestionBanksLogic) successResp(list []*types.QuestionBank, total int64) (*types.GetQuestionBankListResp, error) {
	return &types.GetQuestionBankListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.QuestionBankList{
			Total: total,
			List:  list,
		},
	}, nil
}
