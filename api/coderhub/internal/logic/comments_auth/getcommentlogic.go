package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/client/commentservice"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCommentLogic 获取单个评论
func NewGetCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLogic {
	return &GetCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentLogic) GetComment(req *types.GetCommentReq) (resp *types.GetCommentResp, err error) {
	comment, err := l.svcCtx.CommentService.GetComment(l.ctx, &commentservice.GetCommentRequest{
		CommentId: req.CommentId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(comment)
}

func (l *GetCommentLogic) successResp(comment *commentservice.GetCommentResponse) (*types.GetCommentResp, error) {
	if comment == nil {
		return &types.GetCommentResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpNotFound,
				Message: conf.HttpMessage.MsgFailed,
			},
		}, nil
	}

	// 构建用户信息，处理空值情况
	var userInfo *types.UserInfo
	if comment.Comment.UserInfo != nil {
		userInfo = &types.UserInfo{
			Id:       comment.Comment.UserInfo.UserId,
			Username: comment.Comment.UserInfo.UserName,
			Nickname: comment.Comment.UserInfo.UserName,
			Email:    comment.Comment.UserInfo.Email,
			Phone:    comment.Comment.UserInfo.Phone,
			Avatar:   comment.Comment.UserInfo.Avatar,
			Gender:   comment.Comment.UserInfo.Gender,
			Age:      comment.Comment.UserInfo.Age,
			Status:   comment.Comment.UserInfo.Status,
			IsAdmin:  comment.Comment.UserInfo.IsAdmin,
			CreateAt: comment.Comment.UserInfo.CreatedAt,
			UpdateAt: comment.Comment.UserInfo.UpdatedAt,
		}
	}

	var replyToUserInfo *types.UserInfo
	if comment.Comment.ReplyToUserInfo != nil {
		replyToUserInfo = &types.UserInfo{
			Id:       comment.Comment.ReplyToUserInfo.UserId,
			Username: comment.Comment.ReplyToUserInfo.UserName,
			Avatar:   comment.Comment.ReplyToUserInfo.Avatar,
			Nickname: comment.Comment.ReplyToUserInfo.UserName,
			Email:    comment.Comment.ReplyToUserInfo.Email,
			Phone:    comment.Comment.ReplyToUserInfo.Phone,
			Gender:   comment.Comment.ReplyToUserInfo.Gender,
			Age:      comment.Comment.ReplyToUserInfo.Age,
			Status:   comment.Comment.ReplyToUserInfo.Status,
			IsAdmin:  comment.Comment.ReplyToUserInfo.IsAdmin,
			CreateAt: comment.Comment.ReplyToUserInfo.CreatedAt,
			UpdateAt: comment.Comment.ReplyToUserInfo.UpdatedAt,
		}
	}

	// 获取图片
	images := make([]types.ImageInfo, len(comment.Comment.Images))
	for i, image := range comment.Comment.Images {
		images[i] = types.ImageInfo{
			ImageId:      image.ImageId,
			BucketName:   image.BucketName,
			ObjectName:   image.ObjectName,
			Url:          image.Url,
			ThumbnailUrl: image.ThumbnailUrl,
			ContentType:  image.ContentType,
			Size:         image.Size,
			Width:        image.Width,
			Height:       image.Height,
			UploadIp:     image.UploadIp,
			UserId:       image.UserId,
			CreatedAt:    image.CreatedAt,
		}
	}

	return &types.GetCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.Comment{
			Id:              comment.Comment.Id,
			EntityID:        comment.Comment.EntityId,
			Content:         comment.Comment.Content,
			ParentId:        comment.Comment.ParentId,
			RootId:          comment.Comment.RootId,
			UserInfo:        userInfo,
			CreatedAt:       comment.Comment.CreatedAt,
			UpdatedAt:       comment.Comment.UpdatedAt,
			Replies:         nil,
			RepliesCount:    comment.Comment.RepliesCount,
			ReplyToUserInfo: replyToUserInfo,
			LikeCount:       comment.Comment.LikeCount,
			Images:          images,
		},
	}, nil
}

func (l *GetCommentLogic) errorResp(err error) (*types.GetCommentResp, error) {
	return &types.GetCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
