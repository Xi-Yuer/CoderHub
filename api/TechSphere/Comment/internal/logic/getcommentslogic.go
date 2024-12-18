package logic

import (
	"context"

	"coderhub/api/TechSphere/Comment/internal/svc"
	"coderhub/api/TechSphere/Comment/internal/types"
	"coderhub/conf"
	"coderhub/rpc/TechSphere/Comment/commentservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCommentsLogic 获取评论列表
func NewGetCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentsLogic {
	return &GetCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentsLogic) GetComments(req *types.GetCommentsReq) (resp *types.GetCommentsResp, err error) {
	comments, err := l.svcCtx.CommentService.GetComments(l.ctx, &commentservice.GetCommentsRequest{
		ArticleId: req.ArticleId,
		Page:      req.Page,
		PageSize:  req.PageSize,
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
		images := make([]types.CommentImage, 0, len(val.Images))
		for _, image := range val.Images {
			images = append(images, types.CommentImage{
				ImageId:      image.ImageId,
				Url:          image.Url,
				ThumbnailUrl: image.ThumbnailUrl,
			})
		}

		var replyToUserInfo *types.UserInfo
		if val.ReplyToUserInfo != nil {
			replyToUserInfo = &types.UserInfo{
				UserId:   val.ReplyToUserInfo.UserId,
				Username: val.ReplyToUserInfo.Username,
				Avatar:   val.ReplyToUserInfo.Avatar,
			}
		}

		rootComments[i] = &types.Comment{
			Id:        val.Id,
			ArticleId: val.ArticleId,
			Content:   val.Content,
			ParentId:  val.ParentId,
			RootId:    val.RootId,
			UserInfo: &types.UserInfo{
				UserId:   val.UserInfo.UserId,
				Username: val.UserInfo.Username,
				Avatar:   val.UserInfo.Avatar,
			},
			ReplyToUserInfo: replyToUserInfo,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
			Replies:      l.buildTree(val.Replies),
			RepliesCount: val.RepliesCount,
			LikeCount:    val.LikeCount,
			Images:       images,
		}
	}
	return rootComments
}
