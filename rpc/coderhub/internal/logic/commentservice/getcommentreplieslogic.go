package commentservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"sort"
	"sync"

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
		l.Logger.Errorf("获取回复列表失败: %v", err)
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

	// 使用并发异步请求图片和用户信息
	var wg sync.WaitGroup
	var imageErr, userErr error
	var imageRelations *coderhub.BatchGetImagesByEntityResponse
	var users *coderhub.BatchGetUserByIDResponse

	// 异步获取图片
	wg.Add(1)
	go func() {
		defer wg.Done()
		imageRelations, imageErr = imagerelationservicelogic.NewBatchGetImagesByEntityLogic(l.ctx, l.svcCtx).BatchGetImagesByEntity(&coderhub.BatchGetImagesByEntityRequest{
			EntityIds:  replyIds,
			EntityType: model.ImageRelationComment,
		})
	}()

	// 异步获取用户信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		users, userErr = userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx).BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
			UserIds: allUserIds,
		})
	}()

	// 等待并发请求完成
	wg.Wait()

	// 错误处理
	if imageErr != nil {
		l.Logger.Errorf("获取回复图片失败: %v", imageErr)
		return nil, imageErr
	}
	if userErr != nil {
		l.Logger.Errorf("获取用户信息失败: %v", userErr)
		return nil, userErr
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

	// 构建用户信息映射
	userInfos := make(map[int64]*coderhub.UserInfo)
	for _, user := range users.UserInfos {
		if user != nil {
			userInfos[user.UserId] = user
		}
	}

	// 获取评论点赞数
	likeCountMap, err := l.svcCtx.CommentRelationLikeRepository.BatchList(l.ctx, replyIds)
	if err != nil {
		l.Logger.Errorf("获取评论点赞数失败: %v", err)
		return nil, err
	}

	// 获取评论点赞状态
	likeStatusMap, err := l.svcCtx.CommentRelationLikeRepository.BatchGetCommentsHasBeenUserLiked(l.ctx, replyIds, in.UserId)
	if err != nil {
		l.Logger.Errorf("获取评论点赞状态失败: %v", err)
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
			EntityAuthorId:  reply.EntityAuthorID,
			UserInfo:        userInfos[reply.UserID], // 直接获取映射中的用户信息
			ReplyToUserInfo: userInfos[reply.ReplyToUID],
			CreatedAt:       reply.CreatedAt.Unix(),
			UpdatedAt:       reply.UpdatedAt.Unix(),
			Replies:         nil,
			RepliesCount:    0,
			LikeCount:       int32(likeCountMap[reply.ID]),
			Images:          replyImages[reply.ID],
			IsLiked:         likeStatusMap[reply.ID], // 直接获取映射中的点赞状态
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
