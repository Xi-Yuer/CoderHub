package comments_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/client/commentservice"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取评论列表
func NewGetCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsLogic {
	return &GetCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentsLogic) GetComments(req *types.GetCommentsReq) (resp *types.GetCommentsResp, err error) {
	comments, err := l.svcCtx.CommentService.GetComments(l.ctx, &commentservice.GetCommentsRequest{
		EntityId: req.EntityID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(comments)
}

func (l *GetCommentsLogic) successResp(comments *commentservice.GetCommentsResponse) (*types.GetCommentsResp, error) {
	rootComments := l.buildTree(comments.Comments)
	return &types.GetCommentsResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.List{
			List:  rootComments, // 只返回顶级评论，子评论在replies字段中
			Total: comments.Total,
		},
	}, nil
}

func (l *GetCommentsLogic) errorResp(err error) (*types.GetCommentsResp, error) {
	return &types.GetCommentsResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}

// buildTree 构建树形结构
func (l *GetCommentsLogic) buildTree(comments []*commentservice.Comment) []*types.Comment {
	rootComments := make([]*types.Comment, len(comments))
	for i, val := range comments {
		// 获取图片
		images := make([]types.ImageInfo, 0, len(val.Images))
		for _, image := range val.Images {
			images = append(images, types.ImageInfo{
				ImageId:      image.ImageId,
				ObjectName:   image.ObjectName,
				BucketName:   image.BucketName,
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

		var replyToUserInfo *types.UserInfo
		if val.ReplyToUserInfo != nil {
			replyToUserInfo = &types.UserInfo{
				Id:       val.Id,
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

		rootComments[i] = &types.Comment{
			Id:       val.Id,
			EntityID: val.EntityId,
			Content:  val.Content,
			ParentId: val.ParentId,
			RootId:   val.RootId,
			UserInfo: &types.UserInfo{
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
			},
			ReplyToUserInfo: replyToUserInfo,
			CreatedAt:       val.CreatedAt,
			UpdatedAt:       val.UpdatedAt,
			Replies:         l.buildTree(val.Replies),
			RepliesCount:    val.RepliesCount,
			LikeCount:       val.LikeCount,
			Images:          images,
		}
	}
	return rootComments
}
