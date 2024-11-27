// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: comment.proto

package commentservice

import (
	"context"

	"coderhub/rpc/TechSphere/Comment/comment"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Comment               = comment.Comment
	CommentImage          = comment.CommentImage
	CreateCommentRequest  = comment.CreateCommentRequest
	CreateCommentResponse = comment.CreateCommentResponse
	DeleteCommentRequest  = comment.DeleteCommentRequest
	DeleteCommentResponse = comment.DeleteCommentResponse
	GetCommentRequest     = comment.GetCommentRequest
	GetCommentResponse    = comment.GetCommentResponse
	GetCommentsRequest    = comment.GetCommentsRequest
	GetCommentsResponse   = comment.GetCommentsResponse

	CommentService interface {
		// 创建评论
		CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
		// 获取评论列表
		GetComments(ctx context.Context, in *GetCommentsRequest, opts ...grpc.CallOption) (*GetCommentsResponse, error)
		// 获取单个评论详情
		GetComment(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
		// 删除评论
		DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*DeleteCommentResponse, error)
	}

	defaultCommentService struct {
		cli zrpc.Client
	}
)

func NewCommentService(cli zrpc.Client) CommentService {
	return &defaultCommentService{
		cli: cli,
	}
}

// 创建评论
func (m *defaultCommentService) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	client := comment.NewCommentServiceClient(m.cli.Conn())
	return client.CreateComment(ctx, in, opts...)
}

// 获取评论列表
func (m *defaultCommentService) GetComments(ctx context.Context, in *GetCommentsRequest, opts ...grpc.CallOption) (*GetCommentsResponse, error) {
	client := comment.NewCommentServiceClient(m.cli.Conn())
	return client.GetComments(ctx, in, opts...)
}

// 获取单个评论详情
func (m *defaultCommentService) GetComment(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	client := comment.NewCommentServiceClient(m.cli.Conn())
	return client.GetComment(ctx, in, opts...)
}

// 删除评论
func (m *defaultCommentService) DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*DeleteCommentResponse, error) {
	client := comment.NewCommentServiceClient(m.cli.Conn())
	return client.DeleteComment(ctx, in, opts...)
}