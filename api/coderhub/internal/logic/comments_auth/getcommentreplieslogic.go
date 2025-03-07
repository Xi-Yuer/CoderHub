package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
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
	userID, _ := utils.GetUserID(l.ctx)
	reply, err := l.svcCtx.CommentService.GetCommentReplies(l.ctx, &coderhub.GetCommentRepliesRequest{
		CommentId: utils.String2Int(req.CommentId),
		Page:      req.Page,
		PageSize:  req.PageSize,
		UserId:    userID,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(reply)
}

func (l *GetCommentRepliesLogic) successResp(reply *coderhub.GetCommentRepliesResponse) (*types.GetCommentRepliesResp, error) {
	// 优化后的评论切片初始化
	replies := make([]*types.Comment, len(reply.Replies))

	for i, val := range reply.Replies {
		// 转换图片信息，避免重复转换
		images := make([]types.ImageInfo, len(val.Images))
		for j, img := range val.Images {
			images[j] = types.ImageInfo{
				ImageId:      utils.Int2String(img.ImageId),
				BucketName:   img.BucketName,
				ObjectName:   img.ObjectName,
				Url:          img.Url,
				ThumbnailUrl: img.ThumbnailUrl,
				ContentType:  img.ContentType,
				Size:         img.Size,
				Width:        img.Width,
				Height:       img.Height,
				UploadIp:     img.UploadIp,
				UserId:       utils.Int2String(img.UserId),
				CreatedAt:    img.CreatedAt,
			}
		}
		var replyToUserInfo *types.UserInfo
		if val.ReplyToUserInfo != nil { // 从 val.ReplyToUserInfo 构建回复目标的用户信息
			replyToUserInfo = &types.UserInfo{
				Id:       utils.Int2String(val.ReplyToUserInfo.UserId),
				Username: val.ReplyToUserInfo.UserName,
				Nickname: val.ReplyToUserInfo.NickName,
				Email:    val.ReplyToUserInfo.Email,
				Phone:    val.ReplyToUserInfo.Phone,
				Avatar:   val.ReplyToUserInfo.Avatar,
				Gender:   val.ReplyToUserInfo.Gender,
				Age:      val.ReplyToUserInfo.Age,
				Status:   val.ReplyToUserInfo.Status,
				IsAdmin:  val.ReplyToUserInfo.IsAdmin,
				CreateAt: val.ReplyToUserInfo.CreatedAt,
				UpdateAt: val.ReplyToUserInfo.UpdatedAt,
			}
		}
		replies[i] = &types.Comment{
			Id:             utils.Int2String(val.Id),
			EntityID:       utils.Int2String(val.EntityId),
			Content:        val.Content,
			RootId:         utils.Int2String(val.RootId),
			ParentId:       utils.Int2String(val.ParentId),
			EntityAuthorId: utils.Int2String(val.EntityAuthorId),
			UserInfo: &types.UserInfo{
				Id:       utils.Int2String(val.UserInfo.UserId),
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
			},
			CreatedAt:       val.CreatedAt,
			UpdatedAt:       val.UpdatedAt,
			Replies:         nil,
			ReplyToUserInfo: replyToUserInfo,
			RepliesCount:    val.RepliesCount,
			LikeCount:       val.LikeCount,
			IsLiked:         val.IsLiked,
			Images:          images,
		}
	}

	// 构建并返回响应
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
