package commentservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsLogic {
	return &GetCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetComments 获取评论列表
func (l *GetCommentsLogic) GetComments(in *coderhub.GetCommentsRequest) (*coderhub.GetCommentsResponse, error) {
	comments, total, err := l.svcCtx.CommentRepository.ListByArticleID(l.ctx, in.EntityId, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	rootComments := l.buildTree(comments)

	return &coderhub.GetCommentsResponse{
		Comments: rootComments,
		Total:    int32(total),
	}, nil
}

// buildTree 构建树形结构
func (l *GetCommentsLogic) buildTree(comments []model.Comment) []*coderhub.Comment {
	if len(comments) == 0 {
		return nil
	}

	// 收集所有评论ID
	commentIds := make([]int64, len(comments))
	for i, val := range comments {
		commentIds[i] = val.ID
	}

	l.Logger.Infof("正在获取评论的图片关联，评论IDs: %v", commentIds)

	batchGetImagesByEntityService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	imageRelations, err := batchGetImagesByEntityService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  commentIds,
		EntityType: model.ImageRelationComment,
	})
	if err != nil {
		l.Logger.Errorf("获取评论图片失败: %v", err)
		return make([]*coderhub.Comment, 0)
	}

	// 构建评论ID到图片列表的映射
	commentImages := make(map[int64][]*coderhub.ImageInfo)
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

	// 收集所有评论者和被回复者ID
	userIds := make([]int64, len(comments))
	replyToUserIds := make([]int64, len(comments))
	for i, val := range comments {
		userIds[i] = val.UserID
		replyToUserIds[i] = val.ReplyToUID
	}
	batchGetUserByIDService := userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx)
	users, err := batchGetUserByIDService.BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
		UserIds: userIds,
	})
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return make([]*coderhub.Comment, 0)
	}
	// 构建用户信息映射
	userInfos := make(map[int64]*coderhub.UserInfo)
	if users != nil && len(users.UserInfos) > 0 {
		for _, user := range users.UserInfos {
			if user != nil {
				l.Logger.Infof("映射用户信息: userId=%d, userName=%s", user.UserId, user.UserName)
				// 如果用户信息不存在，则添加到映射中
				if _, ok := userInfos[user.UserId]; !ok {
					userInfos[user.UserId] = &coderhub.UserInfo{
						UserId:    user.UserId,
						UserName:  user.UserName,
						Avatar:    user.Avatar,
						Email:     user.Email,
						Password:  user.Password,
						Gender:    user.Gender,
						Age:       user.Age,
						Phone:     user.Phone,
						NickName:  user.NickName,
						IsAdmin:   user.IsAdmin,
						Status:    user.Status,
						CreatedAt: user.CreatedAt,
						UpdatedAt: user.UpdatedAt,
					}
				}
			}
		}
	}

	rootComments := make([]*coderhub.Comment, len(comments))
	// 获取评论点赞数
	likeCountMap, err := l.svcCtx.CommentRelationLikeRepository.BatchList(l.ctx, commentIds)
	if err != nil {
		l.Logger.Errorf("获取评论点赞数失败: %v", err)
		return make([]*coderhub.Comment, 0)
	}
	for i, val := range comments {
		// 确保每个评论的图片列表都被初始化
		if _, ok := commentImages[val.ID]; !ok {
			commentImages[val.ID] = make([]*coderhub.ImageInfo, 0)
		}

		rootComments[i] = &coderhub.Comment{
			Id:              val.ID,
			EntityId:        val.EntityID,
			Content:         val.Content,
			ParentId:        val.ParentID,
			RootId:          val.RootID,
			UserInfo:        userInfos[val.UserID],
			ReplyToUserInfo: userInfos[val.ReplyToUID],
			CreatedAt:       val.CreatedAt.Unix(),
			UpdatedAt:       val.UpdatedAt.Unix(),
			Replies:         l.buildTree(val.Replies),
			RepliesCount:    val.ReplyCount,
			LikeCount:       int32(likeCountMap[val.ID]),
			Images:          commentImages[val.ID],
		}
	}
	// 按照点赞数量进行排序
	sort.Slice(rootComments, func(i, j int) bool {
		return rootComments[i].LikeCount > rootComments[j].LikeCount
	})
	return rootComments
}
