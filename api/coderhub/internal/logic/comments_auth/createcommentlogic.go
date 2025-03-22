package comments_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/rpc/coderhub/client/commentservice"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"fmt"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateCommentLogic 创建评论
func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CreateCommentResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	comment, err := l.svcCtx.CommentService.CreateComment(utils.SetUserMetaData(l.ctx), &coderhub.CreateCommentRequest{
		EntityId:       utils.String2Int(req.EntityID),
		Content:        req.Content,
		ParentId:       utils.String2Int(req.ParentId),
		RootId:         utils.String2Int(req.RootId),
		UserId:         userID,
		ReplyToUid:     utils.String2Int(req.ReplyToUID),
		ImageIds:       utils.StringArray2Int64Array(req.ImageIds),
		EntityAuthorId: utils.String2Int(req.EntityAuthorID),
	})
	if err != nil {
		return l.errorResp(err)
	}
	err = l.SendMessage(l.ctx, comment, utils.String2Int(req.EntityID), userID, utils.String2Int(req.ReplyToUID), req.Content, req.ParentId)
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp(comment)
}

// SendMessage 发送评论消息
func (l *CreateCommentLogic) SendMessage(ctx context.Context, comment *coderhub.CreateCommentResponse, entityId, userId, ReplyToUid int64, content string, parentId string) error {
	// 发送评论用户的信息
	userInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{
		UserId:        comment.Comment.UserInfo.UserId,
		RequestUserId: 0,
	})
	if err != nil {
		fmt.Printf("获取用户信息失败：%s\n", err.Error())
		return err
	}
	// 获取文章信息
	article, err := l.svcCtx.ArticlesService.GetArticle(ctx, &coderhub.GetArticleRequest{
		Id:     entityId,
		UserId: 0,
	})
	if err == nil {
		// 发送评论通知
		var title string
		var _type string
		if article.Article.Title != "" {
			title = article.Article.Title
			_type = "文章"
		} else {
			title = "沸点"
			_type = ""
		}
		_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
			SenderId:   userId,
			ReceiverId: article.Article.AuthorId,
			Type:       model.MessageComment,
			EntityId:   article.Article.Id,
			Content:    fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 评论了你的%s <a className=\"font-bold inline-block mx-2\" href=\"/post/%s\" target=\"_blank\">《%s》</a>", utils.Int2String(userInfo.UserId), userInfo.UserName, _type, utils.Int2String(article.Article.Id), title),
		})
		if err != nil {
			fmt.Printf("发送评论通知失败：%s\n", err.Error())
			return err
		}
	}

	// 发送回复通知
	if ReplyToUid != 0 {
		comment, err := l.svcCtx.CommentService.GetComment(ctx, &coderhub.GetCommentRequest{
			CommentId: utils.String2Int(parentId),
			UserId:    0,
		})
		if err != nil {
			return err
		}
		_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
			SenderId:   userId,
			ReceiverId: ReplyToUid,
			Type:       model.MessageComment,
			EntityId:   entityId,
			Content:    fmt.Sprintf("用户 <a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 回复了你的评论：%s <br> <span class=\"text-gray-500 mt-4 text-sm inline-block\">%s</span>", utils.Int2String(userInfo.UserId), userInfo.UserName, content, comment.Comment.Content),
		})
		if err != nil {
			fmt.Printf("发送评论通知失败：%s\n", err.Error())
			return err
		}
	}
	return nil
}

func (l *CreateCommentLogic) successResp(comment *commentservice.CreateCommentResponse) (*types.CreateCommentResp, error) {

	var Images []types.ImageInfo
	for _, image := range comment.Comment.Images {
		Images = append(Images, types.ImageInfo{
			ImageId:      utils.Int2String(image.ImageId),
			BucketName:   image.BucketName,
			ObjectName:   image.ObjectName,
			Url:          image.Url,
			ThumbnailUrl: image.ThumbnailUrl,
			ContentType:  image.ContentType,
			Size:         image.Size,
			Width:        image.Width,
			Height:       image.Height,
			UploadIp:     image.UploadIp,
			UserId:       utils.Int2String(image.UserId),
			CreatedAt:    image.CreatedAt,
		})
	}
	var replyToUserInfo *types.UserInfo
	var userInfo *types.UserInfo
	if comment.Comment.ReplyToUserInfo != nil {
		replyToUserInfo = &types.UserInfo{
			Id:          utils.Int2String(comment.Comment.ReplyToUserInfo.UserId),
			Username:    comment.Comment.ReplyToUserInfo.UserName,
			Nickname:    comment.Comment.ReplyToUserInfo.NickName,
			Email:       comment.Comment.ReplyToUserInfo.Email,
			Phone:       comment.Comment.ReplyToUserInfo.Phone,
			Avatar:      comment.Comment.ReplyToUserInfo.Avatar,
			Gender:      comment.Comment.ReplyToUserInfo.Gender,
			Age:         comment.Comment.ReplyToUserInfo.Age,
			Status:      comment.Comment.ReplyToUserInfo.Status,
			IsAdmin:     comment.Comment.ReplyToUserInfo.IsAdmin,
			CreateAt:    comment.Comment.ReplyToUserInfo.CreatedAt,
			UpdateAt:    comment.Comment.ReplyToUserInfo.UpdatedAt,
			FollowCount: comment.Comment.ReplyToUserInfo.FollowCount,
			FansCount:   comment.Comment.ReplyToUserInfo.UserId,
			IsFollowed:  comment.Comment.ReplyToUserInfo.IsFollowed,
		}
	}
	if comment.Comment.UserInfo != nil {
		userInfo = &types.UserInfo{
			Id:          utils.Int2String(comment.Comment.UserInfo.UserId),
			Username:    comment.Comment.UserInfo.UserName,
			Nickname:    comment.Comment.UserInfo.NickName,
			Email:       comment.Comment.UserInfo.Email,
			Phone:       comment.Comment.UserInfo.Phone,
			Avatar:      comment.Comment.UserInfo.Avatar,
			Gender:      comment.Comment.UserInfo.Gender,
			Age:         comment.Comment.UserInfo.Age,
			Status:      comment.Comment.UserInfo.Status,
			IsAdmin:     comment.Comment.UserInfo.IsAdmin,
			CreateAt:    comment.Comment.UserInfo.CreatedAt,
			UpdateAt:    comment.Comment.UserInfo.UpdatedAt,
			FollowCount: comment.Comment.UserInfo.FollowerCount,
			FansCount:   comment.Comment.UserInfo.FollowCount,
			IsFollowed:  comment.Comment.UserInfo.IsFollowed,
		}
	}
	return &types.CreateCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.Comment{
			Id:              utils.Int2String(comment.Comment.Id),
			EntityID:        utils.Int2String(comment.Comment.EntityId),
			Content:         comment.Comment.Content,
			RootId:          utils.Int2String(comment.Comment.RootId),
			ParentId:        utils.Int2String(comment.Comment.ParentId),
			EntityAuthorId:  utils.Int2String(comment.Comment.EntityAuthorId),
			UserInfo:        userInfo,
			CreatedAt:       comment.Comment.CreatedAt,
			UpdatedAt:       comment.Comment.UpdatedAt,
			Replies:         nil,
			ReplyToUserInfo: replyToUserInfo,
			RepliesCount:    comment.Comment.RepliesCount,
			LikeCount:       comment.Comment.LikeCount,
			Images:          Images,
		},
	}, nil
}

func (l *CreateCommentLogic) errorResp(err error) (*types.CreateCommentResp, error) {
	return &types.CreateCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}
