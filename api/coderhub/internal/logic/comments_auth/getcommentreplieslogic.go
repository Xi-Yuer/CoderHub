package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentRepliesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCommentRepliesLogic 获取某条评论的子评论列表
func NewGetCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentRepliesLogic {
	return &GetCommentRepliesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentRepliesLogic) GetCommentReplies(req *types.GetCommentRepliesReq) (resp *types.GetCommentRepliesResp, err error) {
	reply, err := l.svcCtx.CommentService.GetCommentReplies(l.ctx, &coderhub.GetCommentRepliesRequest{
		CommentId: req.CommentId,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(reply)
}

func (l *GetCommentRepliesLogic) successResp(reply *coderhub.GetCommentRepliesResponse) (*types.GetCommentRepliesResp, error) {
	replies := make([]*types.Comment, len(reply.Replies))
	for i, val := range reply.Replies {
		images := make([]types.ImageInfo, len(val.Images))
		for j, img := range val.Images {
			images[j] = types.ImageInfo{
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
			}
		}
		var replyToUserInfo *types.UserInfo
		if val.ReplyToUserInfo != nil {
			replyToUserInfo = &types.UserInfo{
				Id:       val.Id,
				Username: val.UserInfo.UserName,
				Nickname: val.UserInfo.NickName,
				Email:    val.UserInfo.Email,
				Phone:    val.UserInfo.Phone,
				Avatar:   val.UserInfo.Avatar,
				Gender:   val.UserInfo.Gender,
				Age:      val.UserInfo.Age,
				Status:   val.UserInfo.Status,
				IsAdmin:  val.UserInfo.IsAdmin,
				CreateAt: val.UserInfo.CreatedAt,
				UpdateAt: val.UserInfo.UpdatedAt,
			}
		}
		replies[i] = &types.Comment{
			Id:              val.Id,
			EntityID:        val.EntityId,
			Content:         val.Content,
			ParentId:        val.ParentId,
			RootId:          val.RootId,
			UserInfo:        replyToUserInfo,
			ReplyToUserInfo: replyToUserInfo,
			CreatedAt:       val.CreatedAt,
			UpdatedAt:       val.UpdatedAt,
			Replies:         nil,
			RepliesCount:    val.RepliesCount,
			LikeCount:       val.LikeCount,
			Images:          images,
		}
	}
	return &types.GetCommentRepliesResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.List{
			List:  replies,
			Total: reply.Total,
		},
	}, nil
}

func (l *GetCommentRepliesLogic) errorResp(err error) (*types.GetCommentRepliesResp, error) {
	return &types.GetCommentRepliesResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: conf.HttpMessage.MsgFailed,
		},
	}, err
}
