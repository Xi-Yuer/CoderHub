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

type GetQuestionBankDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetQuestionBankDetailLogic 获取题库详情
func NewGetQuestionBankDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionBankDetailLogic {
	return &GetQuestionBankDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionBankDetailLogic) GetQuestionBankDetail(req *types.GetQuestionBankDetailReq) (resp *types.GetQuestionBankDetailResp, err error) {
	l.Logger.Info("GetQuestionBankDetail", req.Id)
	l.Logger.Info("GetQuestionBankDetail", req.RequestUserID)
	questionBank, err := l.svcCtx.QuestionBankService.GetQuestionBank(l.ctx, &coderhub.GetQuestionBankRequest{
		BankId: utils.String2Int(req.Id),
		UserId: utils.String2Int(req.RequestUserID),
	})
	if err != nil {
		return l.error(err)
	}
	return l.success(questionBank)
}

func (l *GetQuestionBankDetailLogic) success(questionBank *coderhub.GetQuestionBankResponse) (resp *types.GetQuestionBankDetailResp, err error) {
	return &types.GetQuestionBankDetailResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.GetQuestionBankResqonse{
			QuestionBank: types.QuestionBank{
				Id:          utils.Int2String(questionBank.Bank.Id),
				Name:        questionBank.Bank.Name,
				Description: questionBank.Bank.Description,
				Difficulty:  questionBank.Bank.Difficulty,
				Tags:        questionBank.Bank.Tags,
				CoverImage: &types.ImageInfo{
					ImageId:      utils.Int2String(questionBank.Bank.CoverImage.ImageId),
					BucketName:   questionBank.Bank.CoverImage.BucketName,
					ObjectName:   questionBank.Bank.CoverImage.ObjectName,
					Url:          questionBank.Bank.CoverImage.Url,
					ThumbnailUrl: questionBank.Bank.CoverImage.ThumbnailUrl,
					ContentType:  questionBank.Bank.CoverImage.ContentType,
					Size:         questionBank.Bank.CoverImage.Size,
					Width:        questionBank.Bank.CoverImage.Width,
					Height:       questionBank.Bank.CoverImage.Height,
					UploadIp:     questionBank.Bank.CoverImage.UploadIp,
					UserId:       utils.Int2String(questionBank.Bank.CoverImage.UserId),
					CreatedAt:    questionBank.Bank.CoverImage.CreatedAt,
				},
				CreateUser: &types.UserInfo{
					Id:          utils.Int2String(questionBank.Bank.CreateUser.UserId),
					Username:    questionBank.Bank.CreateUser.UserName,
					Nickname:    questionBank.Bank.CreateUser.NickName,
					Email:       questionBank.Bank.CreateUser.Email,
					Phone:       questionBank.Bank.CreateUser.Phone,
					Avatar:      questionBank.Bank.CreateUser.Avatar,
					Gender:      questionBank.Bank.CreateUser.Gender,
					Age:         questionBank.Bank.CreateUser.Age,
					Status:      questionBank.Bank.CreateUser.Status,
					IsAdmin:     questionBank.Bank.CreateUser.IsAdmin,
					CreateAt:    questionBank.Bank.CreateUser.CreatedAt,
					UpdateAt:    questionBank.Bank.CreateUser.UpdatedAt,
					FollowCount: questionBank.Bank.CreateUser.FollowerCount,
					FansCount:   questionBank.Bank.CreateUser.FollowCount,
					IsFollowed:  questionBank.Bank.CreateUser.IsFollowed,
				},
				CreatedAt: questionBank.Bank.CreateTime,
				UpdatedAt: questionBank.Bank.UpdateTime,
			},
			IsFavorited: questionBank.IsFavorite,
		},
	}, nil
}

func (l *GetQuestionBankDetailLogic) error(error error) (resp *types.GetQuestionBankDetailResp, err error) {
	return &types.GetQuestionBankDetailResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: error.Error(),
		},
		Data: nil,
	}, nil
}
