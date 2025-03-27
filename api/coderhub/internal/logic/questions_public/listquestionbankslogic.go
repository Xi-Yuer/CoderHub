package questions_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
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
		CategoryId: utils.String2Int(req.CategoryID),
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	// 确保 list 不是 nil
	if list == nil || list.Banks == nil {
		return l.successResp(nil, 0)
	}

	banks := make([]*types.QuestionBank, 0, len(list.Banks))
	for _, bank := range list.Banks {
		// 确保 CoverImage 不是 nil
		var coverImage *types.ImageInfo
		if bank.CoverImage != nil {
			coverImage = &types.ImageInfo{
				ImageId:      utils.Int2String(bank.CoverImage.ImageId),
				BucketName:   bank.CoverImage.BucketName,
				ObjectName:   bank.CoverImage.ObjectName,
				Url:          bank.CoverImage.Url,
				ThumbnailUrl: bank.CoverImage.ThumbnailUrl,
				ContentType:  bank.CoverImage.ContentType,
				Size:         bank.CoverImage.Size,
				Width:        bank.CoverImage.Width,
				Height:       bank.CoverImage.Height,
				UploadIp:     bank.CoverImage.UploadIp,
				UserId:       utils.Int2String(bank.CoverImage.UserId),
				CreatedAt:    bank.CoverImage.CreatedAt,
			}
		}

		// 确保 CreateUser 不是 nil
		var createUser *types.UserInfo
		if bank.CreateUser != nil {
			createUser = &types.UserInfo{
				Id:       utils.Int2String(bank.CreateUser.UserId),
				Username: bank.CreateUser.UserName,
				Nickname: bank.CreateUser.NickName,
				Email:    bank.CreateUser.Email,
				Phone:    bank.CreateUser.Phone,
				Avatar:   bank.CreateUser.Avatar,
				Gender:   bank.CreateUser.Gender,
				Age:      bank.CreateUser.Age,
				Status:   bank.CreateUser.Status,
				IsAdmin:  bank.CreateUser.IsAdmin,
				CreateAt: bank.CreateUser.CreatedAt,
				UpdateAt: bank.CreateUser.UpdatedAt,
			}
		}

		banks = append(banks, &types.QuestionBank{
			Id:          utils.Int2String(bank.Id),
			Name:        bank.Name,
			Description: bank.Description,
			Difficulty:  bank.Difficulty,
			Tags:        bank.Tags,
			CoverImage:  coverImage, // 可能为 nil，但不会 panic
			CreateUser:  createUser, // 可能为 nil，但不会 panic
			CreatedAt:   bank.CreateTime,
			UpdatedAt:   bank.UpdateTime,
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
