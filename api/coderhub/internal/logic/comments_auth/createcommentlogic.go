package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/client/commentservice"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

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
		EntityAuthorId: utils.String2Int(req.EntityAuthorID),
		ReplyToUid:     utils.String2Int(req.ReplyToUID),
		ImageIds:       utils.StringArray2Int64Array(req.ImageIds),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(comment)
}

func (l *CreateCommentLogic) successResp(comment *commentservice.CreateCommentResponse) (*types.CreateCommentResp, error) {

	Images := make([]types.ImageInfo, len(comment.Comment.Images))
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
			Id:       utils.Int2String(comment.Comment.ReplyToUserInfo.UserId),
			Username: comment.Comment.ReplyToUserInfo.UserName,
			Nickname: comment.Comment.ReplyToUserInfo.NickName,
			Email:    comment.Comment.ReplyToUserInfo.Email,
			Phone:    comment.Comment.ReplyToUserInfo.Phone,
			Avatar:   comment.Comment.ReplyToUserInfo.Avatar,
		}
	}
	if comment.Comment.UserInfo != nil {
		userInfo = &types.UserInfo{
			Id:       utils.Int2String(comment.Comment.UserInfo.UserId),
			Username: comment.Comment.UserInfo.UserName,
			Nickname: comment.Comment.UserInfo.NickName,
			Email:    comment.Comment.UserInfo.Email,
			Phone:    comment.Comment.UserInfo.Phone,
			Avatar:   comment.Comment.UserInfo.Avatar,
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
