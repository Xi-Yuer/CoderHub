package questionservicelogic

import (
	"coderhub/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionBankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionBankLogic {
	return &GetQuestionBankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetQuestionBank 获取题库详情
// GetQuestionBank 获取题库详情
func (l *GetQuestionBankLogic) GetQuestionBank(in *coderhub.GetQuestionBankRequest) (*coderhub.GetQuestionBankResponse, error) {
	// 获取题库信息
	questionBank, err := l.svcCtx.QuestionBankRepository.GetQuestionBankByID(l.ctx, in.BankId)
	if err != nil {
		return nil, fmt.Errorf("查询题库失败: %w", err)
	}
	if questionBank == nil {
		return nil, errors.New("题库不存在")
	}

	// 解析 Tags
	var tags []string
	if questionBank.Tags != "" {
		tags = strings.Split(questionBank.Tags, ",")
	}

	// 获取题库封面图片
	var CoverImage *coderhub.ImageInfo
	ids, err := l.svcCtx.ImageRelationRepository.ListByEntityID(l.ctx, questionBank.ID, model.ImageRelationQuestionCover)
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		image, err := l.svcCtx.ImageRepository.GetByID(l.ctx, ids[0].ImageID)
		if err != nil {
			return nil, fmt.Errorf("查询封面图片失败: %w", err)
		}
		if image != nil {
			CoverImage = &coderhub.ImageInfo{
				ImageId:      image.ID,
				BucketName:   image.BucketName,
				ObjectName:   image.ObjectName,
				Url:          image.URL,
				ThumbnailUrl: image.ThumbnailURL,
				ContentType:  image.ContentType,
				Size:         image.Size,
				Width:        image.Width,
				Height:       image.Height,
				UploadIp:     image.UploadIP,
				UserId:       image.UserID,
				CreatedAt:    image.CreatedAt.Unix(),
			}
		}
	}

	// 获取创建用户信息（优化为单次查询）
	user, err := l.svcCtx.UserRepository.GetUserByID(questionBank.CreateUser)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("创建用户不存在")
	}

	// 判断是否被收藏（避免无效查询）
	exist := false
	if user.ID > 0 && in.BankId > 0 {
		exist, err = l.svcCtx.UserFavorEntityRepository.IsFavorEntityExist(l.ctx, in.UserId, in.BankId, model.UserFavorRelationQuestion)
		if err != nil {
			return nil, fmt.Errorf("查询收藏状态失败: %w", err)
		}
	}

	// 返回优化后的响应
	return &coderhub.GetQuestionBankResponse{
		Bank: &coderhub.QuestionBank{
			Id:          questionBank.ID,
			Description: questionBank.Description,
			Difficulty:  questionBank.Difficulty,
			Tags:        tags,
			CreateUser: &coderhub.UserInfo{
				UserId:        user.ID,
				UserName:      user.UserName,
				Avatar:        user.Avatar.String,
				Email:         user.Email.String,
				Password:      user.Password,
				Gender:        user.Gender,
				Age:           user.Age,
				Phone:         user.Phone.String,
				NickName:      user.NickName.String,
				IsAdmin:       user.IsAdmin,
				Status:        user.Status,
				CreatedAt:     user.CreatedAt.Unix(),
				UpdatedAt:     user.UpdatedAt.Unix(),
				FollowCount:   user.FollowerCount,
				FollowerCount: user.FollowCount,
				IsFollowed:    user.IsFollowed,
			},
			Name:       questionBank.Name,
			CoverImage: CoverImage,
			CreateTime: questionBank.CreatedAt.Unix(),
			UpdateTime: questionBank.UpdatedAt.Unix(),
		},
		IsFavorite: exist,
	}, nil
}
