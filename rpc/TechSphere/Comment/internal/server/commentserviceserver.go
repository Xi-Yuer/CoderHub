// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: comment.proto

package server

import (
	"context"

	"coderhub/rpc/TechSphere/Comment/comment"
	"coderhub/rpc/TechSphere/Comment/internal/logic"
	"coderhub/rpc/TechSphere/Comment/internal/svc"
)

type CommentServiceServer struct {
	svcCtx *svc.ServiceContext
	comment.UnimplementedCommentServiceServer
}

func NewCommentServiceServer(svcCtx *svc.ServiceContext) *CommentServiceServer {
	return &CommentServiceServer{
		svcCtx: svcCtx,
	}
}

// 创建评论
func (s *CommentServiceServer) CreateComment(ctx context.Context, in *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	l := logic.NewCreateCommentLogic(ctx, s.svcCtx)
	return l.CreateComment(in)
}

// 获取评论列表
func (s *CommentServiceServer) GetComments(ctx context.Context, in *comment.GetCommentsRequest) (*comment.GetCommentsResponse, error) {
	l := logic.NewGetCommentsLogic(ctx, s.svcCtx)
	return l.GetComments(in)
}

// 获取单个评论详情
func (s *CommentServiceServer) GetComment(ctx context.Context, in *comment.GetCommentRequest) (*comment.GetCommentResponse, error) {
	l := logic.NewGetCommentLogic(ctx, s.svcCtx)
	return l.GetComment(in)
}

// 更新评论
func (s *CommentServiceServer) UpdateComment(ctx context.Context, in *comment.UpdateCommentRequest) (*comment.UpdateCommentResponse, error) {
	l := logic.NewUpdateCommentLogic(ctx, s.svcCtx)
	return l.UpdateComment(in)
}

// 删除评论
func (s *CommentServiceServer) DeleteComment(ctx context.Context, in *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	l := logic.NewDeleteCommentLogic(ctx, s.svcCtx)
	return l.DeleteComment(in)
}
