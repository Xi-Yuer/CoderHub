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
	ctx := utils.SetUserMetaData(l.ctx)
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	comment, err := l.svcCtx.CommentService.CreateComment(ctx, &coderhub.CreateCommentRequest{
		EntityId:   req.EntityID,
		Content:    req.Content,
		ParentId:   req.ParentId,
		RootId:     req.RootId,
		UserId:     userID,
		ReplyToUid: req.ReplyToUID,
		ImageIds:   req.ImageIds,
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
		})
	}
	return &types.CreateCommentResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.Comment{
			Id:           comment.Comment.Id,
			EntityID:     comment.Comment.EntityId,
			Content:      comment.Comment.Content,
			ParentId:     comment.Comment.ParentId,
			RootId:       comment.Comment.RootId,
			UserInfo:     &types.UserInfo{Id: comment.Comment.UserInfo.UserId, Username: comment.Comment.UserInfo.UserName, Avatar: comment.Comment.UserInfo.Avatar},
			CreatedAt:    comment.Comment.CreatedAt,
			UpdatedAt:    comment.Comment.UpdatedAt,
			Replies:      nil,
			RepliesCount: comment.Comment.RepliesCount,
			LikeCount:    comment.Comment.LikeCount,
			Images:       Images,
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
