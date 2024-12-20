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

type GetCommentRepliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentRepliesLogic {
	return &GetCommentRepliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCommentReplies 获取某条评论的子评论列表
func (l *GetCommentRepliesLogic) GetCommentReplies(in *coderhub.GetCommentRepliesRequest) (*coderhub.GetCommentRepliesResponse, error) {
	// 获取回复列表
	replies, total, err := l.svcCtx.CommentRepository.ListReplies(l.ctx, in.CommentId, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	// 如果没有回复，直接返回空结果
	if len(replies) == 0 {
		return &coderhub.GetCommentRepliesResponse{
			Replies: make([]*coderhub.Comment, 0),
			Total:   0,
		}, nil
	}

	// 收集所有回复ID和用户ID
	replyIds := make([]int64, len(replies))
	userIds := make([]int64, len(replies))
	replyToUserIds := make([]int64, 0)
	for i, reply := range replies {
		replyIds[i] = reply.ID
		userIds[i] = reply.UserID
		if reply.ReplyToUID > 0 {
			replyToUserIds = append(replyToUserIds, reply.ReplyToUID)
		}
	}

	// 合并所有需要查询的用户ID
	allUserIds := append(userIds, replyToUserIds...)

	// 获取图片关联
	batchGetImageByEntityService := imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx)
	imageRelations, err := batchGetImageByEntityService.BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
		EntityIds:  replyIds,
		EntityType: model.ImageRelationComment,
	})
	if err != nil {
		l.Logger.Errorf("获取回复图片失败: %v", err)
		return nil, err
	}

	// 构建回复ID到图片列表的映射
	replyImages := make(map[int64][]*coderhub.ImageInfo)
	for _, img := range imageRelations.Relations {
		if img.ImageId > 0 {
			replyImages[img.EntityId] = append(replyImages[img.EntityId], &coderhub.ImageInfo{
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

	// 获取所有用户信息（包括评论者和被回复者）
	getUserInfoService := userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx)
	users, err := getUserInfoService.BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
		UserIds: allUserIds,
	})
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	// 添加调试日志
	l.Logger.Infof("用户ID列表: %v", allUserIds)
	l.Logger.Infof("获取到的用户信息: %+v", users)

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

	// 获取评论点赞数
	likeCountMap, err := l.svcCtx.CommentRelationLikeRepository.BatchList(l.ctx, replyIds)
	if err != nil {
		l.Logger.Errorf("获取评论点赞数失败: %v", err)
		return nil, err
	}

	// 构建回复列表时添加被回复者信息
	commentReplies := make([]*coderhub.Comment, len(replies))
	for i, reply := range replies {
		if _, ok := replyImages[reply.ID]; !ok {
			replyImages[reply.ID] = make([]*coderhub.ImageInfo, 0)
		}
		commentReplies[i] = &coderhub.Comment{
			Id:              reply.ID,
			EntityId:        reply.EntityID,
			Content:         reply.Content,
			ParentId:        reply.ParentID,
			RootId:          reply.RootID,
			UserInfo:        userInfos[reply.UserID],
			ReplyToUserInfo: userInfos[reply.ReplyToUID],
			LikeCount:       int32(likeCountMap[reply.ID]),
			Images:          replyImages[reply.ID],
			CreatedAt:       reply.CreatedAt.Unix(),
			UpdatedAt:       reply.UpdatedAt.Unix(),
		}
	}

	// 按照点赞数量进行排序
	sort.Slice(commentReplies, func(i, j int) bool {
		return commentReplies[i].LikeCount > commentReplies[j].LikeCount
	})

	return &coderhub.GetCommentRepliesResponse{
		Replies: commentReplies,
		Total:   int32(total),
	}, nil
}
