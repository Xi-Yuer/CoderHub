package comments_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strings"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLikeCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateCommentLikeCountLogic 更新评论点赞数
func NewUpdateCommentLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLikeCountLogic {
	return &UpdateCommentLikeCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentLikeCountLogic) UpdateCommentLikeCount(req *types.UpdateCommentLikeCountReq) (resp *types.UpdateCommentLikeCountResp, err error) {
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err), nil
	}

	if _, err := l.svcCtx.CommentService.UpdateCommentLikeCount(utils.SetUserMetaData(l.ctx), &coderhub.UpdateCommentLikeCountRequest{
		CommentId: utils.String2Int(req.CommentId),
		UserId:    userId,
	}); err != nil {
		return l.errorResp(err), nil
	}

	// 发送消息
	if req.Trigger {
		err = l.SendMessage(l.ctx, utils.String2Int(req.CommentId), userId)
		if err != nil {
			l.errorResp(err)
		}
	}
	return l.successResp(), nil
}

// SendMessage 发送评论点赞消息
func (l *UpdateCommentLikeCountLogic) SendMessage(ctx context.Context, entityId, userId int64) error {
	comment, err := l.svcCtx.CommentService.GetComment(ctx, &coderhub.GetCommentRequest{
		CommentId: entityId,
		UserId:    0,
	})
	if err != nil {
		fmt.Printf("获取评论失败: %v", err)
		return err
	}

	userInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{
		UserId:        comment.Comment.UserInfo.UserId,
		RequestUserId: 0,
	})
	if err != nil {
		fmt.Printf("获取用户信息失败：%s\n", err.Error())
		return err
	}
	commentAbstract := strings.Split(comment.Comment.Content, "\n")
	_, err = l.svcCtx.MessageService.CreateMessage(ctx, &coderhub.CreateMessageRequest{
		SenderId:   userId,
		ReceiverId: comment.Comment.UserInfo.UserId,
		Type:       model.MessageLike,
		EntityId:   comment.Comment.EntityId,
		Content:    fmt.Sprintf("用户<a className=\"font-bold inline-block mx-2\" href=\"/user/%s\" target=\"_blank\">%s</a> 点赞了你的评论 %s", utils.Int2String(userInfo.UserId), userInfo.UserName, commentAbstract),
	})
	if err != nil {
		fmt.Printf("发送评论通知失败：%s\n", err.Error())
		return err
	}
	return nil
}

func (l *UpdateCommentLikeCountLogic) errorResp(err error) *types.UpdateCommentLikeCountResp {
	return &types.UpdateCommentLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}

func (l *UpdateCommentLikeCountLogic) successResp() *types.UpdateCommentLikeCountResp {
	return &types.UpdateCommentLikeCountResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "success",
		},
		Data: true,
	}
}
