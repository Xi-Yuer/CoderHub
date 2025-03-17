package questionservicelogic

import (
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	"context"
	"strings"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionBankListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionBankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionBankListLogic {
	return &GetQuestionBankListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetQuestionBankList 获取题库列表
func (l *GetQuestionBankListLogic) GetQuestionBankList(in *coderhub.GetQuestionBankListRequest) (*coderhub.GetQuestionBankListResponse, error) {
	banks, total, err := l.svcCtx.QuestionBankRepository.GetQuestionBanks(l.ctx, in.CategoryId, in.Page, in.PageSize)

	// 获取封面
	entityIds := make([]int64, 0)
	userIds := make([]int64, 0)
	for _, bank := range banks {
		entityIds = append(entityIds, bank.ID)
		userIds = append(userIds, bank.CreateUser)
	}
	batchGetImagesByEntityService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	imageRelations, err := batchGetImagesByEntityService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  entityIds,
		EntityType: model.ImageRelationQuestionCover,
	})
	if err != nil {
		return nil, err
	}

	l.Logger.Infof("获取图片列表: %v", imageRelations.Relations)

	// 构建评论ID到图片列表的映射
	commentImages := make(map[int64][]*coderhub.ImageInfo, len(imageRelations.Relations))
	for _, img := range imageRelations.Relations {
		l.Logger.Infof("处理图片关联: EntityId=%d, ImageId=%d", img.EntityId, img.ImageId)
		// 只有当图片ID大于0时才处理
		if img.ImageId > 0 {
			commentImages[img.EntityId] = append(commentImages[img.EntityId], &coderhub.ImageInfo{
				ImageId:      img.ImageId,
				BucketName:   img.BucketName,
				ObjectName:   img.ObjectName,
				Url:          img.Url,
				ThumbnailUrl: img.ThumbnailUrl,
				ContentType:  img.ContentType,
				Size:         img.Size,
				Width:        img.Width,
				Height:       img.Height,
				UploadIp:     img.UploadIp,
				UserId:       img.UserId,
				CreatedAt:    img.CreatedAt,
			})
		}
	}

	l.Logger.Infof("构建评论ID到图片列表的映射: %v", commentImages)

	users, err := l.svcCtx.UserRepository.BatchGetUserByID(userIds)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]*coderhub.UserInfo, len(users))
	for _, user := range users {
		userMap[user.ID] = &coderhub.UserInfo{
			UserId:    user.ID,
			UserName:  user.UserName,
			Avatar:    user.Avatar.String,
			Email:     user.Email.String,
			Password:  "",
			Gender:    user.Gender,
			Age:       user.Age,
			Phone:     user.Phone.String,
			NickName:  user.NickName.String,
			IsAdmin:   user.IsAdmin,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		}
	}

	list := make([]*coderhub.QuestionBank, len(banks))
	for i, bank := range banks {
		coverImage := &coderhub.ImageInfo{}
		if images, ok := commentImages[bank.ID]; ok && len(images) > 0 {
			coverImage = &coderhub.ImageInfo{
				ImageId:      images[0].ImageId,
				BucketName:   images[0].BucketName,
				ObjectName:   images[0].ObjectName,
				Url:          images[0].Url,
				ThumbnailUrl: images[0].ThumbnailUrl,
				ContentType:  images[0].ContentType,
				Size:         images[0].Size,
				Width:        images[0].Width,
				Height:       images[0].Height,
				UploadIp:     images[0].UploadIp,
				UserId:       images[0].UserId,
				CreatedAt:    images[0].CreatedAt,
			}
		}
		list[i] = &coderhub.QuestionBank{
			Id:          bank.ID,
			Description: bank.Description,
			Difficulty:  bank.Difficulty,
			Tags:        strings.Split(bank.Tags, ","),
			CreateUser:  userMap[bank.CreateUser],
			Name:        bank.Name,
			CoverImage:  coverImage,
			CreateTime:  bank.CreatedAt.Unix(),
			UpdateTime:  bank.UpdatedAt.Unix(),
		}
	}

	return &coderhub.GetQuestionBankListResponse{
		Banks: list,
		Total: total,
	}, nil
}
