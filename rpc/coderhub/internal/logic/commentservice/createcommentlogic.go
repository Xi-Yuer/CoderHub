package commentservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"coderhub/rpc/coderhub/internal/svc"
	"coderhub/shared/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateComment 创建评论
func (l *CreateCommentLogic) CreateComment(in *coderhub.CreateCommentRequest) (*coderhub.CreateCommentResponse, error) {
	CommentID := utils.GenID()

	// 构建评论模型
	commentModel := &model.Comment{
		ID:             CommentID,
		EntityID:       in.EntityId,
		Content:        in.Content,
		ParentID:       in.ParentId,
		RootID:         in.RootId,
		UserID:         in.UserId,
		ReplyToUID:     in.ReplyToUid,
		EntityAuthorID: in.EntityAuthorId,
	}

	// 获取用户信息
	userService := userservicelogic.NewGetUserInfoLogic(l.ctx, l.svcCtx)
	user, err := userService.GetUserInfo(&coderhub.GetUserInfoRequest{UserId: in.UserId})
	if err != nil {
		return nil, err
	}

	// 处理图片关联
	imageRelationModels := make([]*coderhub.CreateRelationRequest, len(in.ImageIds))
	for i, imageId := range in.ImageIds {
		imageRelationModels[i] = &coderhub.CreateRelationRequest{
			ImageId:    imageId,
			EntityId:   CommentID,
			EntityType: model.ImageRelationComment,
		}
	}

	// 批量创建图片关系
	if len(imageRelationModels) > 0 {
		imageBatchCreateService := imagerelationservicelogic.NewBatchCreateRelationLogic(l.ctx, l.svcCtx)
		if _, err = imageBatchCreateService.BatchCreateRelation(&coderhub.BatchCreateRelationRequest{Relations: imageRelationModels}); err != nil {
			return nil, err
		}
	}

	// 创建评论
	if err = l.svcCtx.CommentRepository.Create(l.ctx, commentModel); err != nil {
		// 事务回滚：删除已创建的图片关联
		if len(imageRelationModels) > 0 {
			imageBatchDeleteService := imagerelationservicelogic.NewBatchDeleteRelationLogic(l.ctx, l.svcCtx)
			_, _ = imageBatchDeleteService.BatchDeleteRelation(&coderhub.BatchDeleteRelationRequest{Ids: []int64{CommentID}})
		}
		return nil, err
	}

	// 批量获取图片信息
	logx.Infof("imageRelationModels: %+v", imageRelationModels)
	var imagesModel = make([]*coderhub.ImageInfo, len(imageRelationModels))
	if len(imageRelationModels) > 0 {
		imageIDs := make([]int64, len(imageRelationModels))
		for i, img := range imageRelationModels {
			imageIDs[i] = img.ImageId
		}

		images, err := l.svcCtx.ImageRepository.BatchGetImagesByID(l.ctx, imageIDs)
		if err != nil {
			return nil, err
		}

		for i, image := range images {
			imagesModel[i] = &coderhub.ImageInfo{
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

	// 返回评论响应
	return &coderhub.CreateCommentResponse{
		Comment: &coderhub.Comment{
			Id:             commentModel.ID,
			EntityId:       commentModel.EntityID,
			Content:        commentModel.Content,
			ParentId:       commentModel.ParentID,
			RootId:         commentModel.RootID,
			UserInfo:       user,
			EntityAuthorId: commentModel.EntityAuthorID,
			CreatedAt:      commentModel.CreatedAt.Unix(),
			UpdatedAt:      commentModel.UpdatedAt.Unix(),
			Replies:        nil,
			RepliesCount:   0,
			LikeCount:      0,
			Images:         imagesModel,
		},
	}, nil
}
