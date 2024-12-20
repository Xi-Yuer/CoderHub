package commentservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLogic {
	return &GetCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetComment 获取单个评论详情
func (l *GetCommentLogic) GetComment(in *coderhub.GetCommentRequest) (*coderhub.GetCommentResponse, error) {
	commentModel, err := l.svcCtx.CommentRepository.GetByID(l.ctx, in.CommentId)
	if err != nil {
		return nil, err
	}
	// 评论ID
	commentId := []int64{commentModel.ID}
	// 获取与评论关联的所有图片
	batchGetImagesByEntityService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	imageRelations, err := batchGetImagesByEntityService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  commentId,
		EntityType: model.ImageRelationComment,
	})
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	getUserInfoService := userservicelogic.NewGetUserInfoLogic(l.ctx, l.svcCtx)
	user, err := getUserInfoService.GetUserInfo(&coderhub.GetUserInfoRequest{
		UserId: commentModel.UserID,
	})
	if err != nil {
		return nil, err
	}

	// 获取被回复者的信息
	replyUserInfo := &coderhub.UserInfo{}
	if commentModel.ReplyToUID != 0 {
		replyUserInfo, err = getUserInfoService.GetUserInfo(&coderhub.GetUserInfoRequest{
			UserId: commentModel.ReplyToUID,
		})
		if err != nil {
			return nil, err
		}
	}

	images := make([]*coderhub.ImageInfo, 0)
	// 将图片关联转换为评论图片
	for _, val := range imageRelations.Relations {
		images = append(images, &coderhub.ImageInfo{
			ImageId:      val.ImageId,
			Url:          val.Url,
			ThumbnailUrl: val.ThumbnailUrl,
		})
	}
	// 获取评论点赞数
	likeCount, err := l.svcCtx.CommentRelationLikeRepository.List(l.ctx, commentModel.ID)
	if err != nil {
		return nil, err
	}
	return &coderhub.GetCommentResponse{
		Comment: &coderhub.Comment{
			Id:              commentModel.ID,
			EntityId:        commentModel.EntityID,
			Content:         commentModel.Content,
			ParentId:        commentModel.ParentID,
			RootId:          commentModel.RootID,
			UserInfo:        user,
			ReplyToUserInfo: replyUserInfo,
			CreatedAt:       commentModel.CreatedAt.Unix(),
			UpdatedAt:       commentModel.UpdatedAt.Unix(),
			Replies:         nil,
			RepliesCount:    0,
			LikeCount:       int32(likeCount),
			Images:          images,
		},
	}, nil
}
