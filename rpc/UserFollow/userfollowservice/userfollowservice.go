// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: user_follow.proto

package userfollowservice

import (
	"context"

	"coderhub/rpc/UserFollow/user_follow"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateUserFollowReq  = user_follow.CreateUserFollowReq
	CreateUserFollowResp = user_follow.CreateUserFollowResp
	DeleteUserFollowReq  = user_follow.DeleteUserFollowReq
	DeleteUserFollowResp = user_follow.DeleteUserFollowResp
	GetMutualFollowsReq  = user_follow.GetMutualFollowsReq
	GetMutualFollowsResp = user_follow.GetMutualFollowsResp
	GetUserFansReq       = user_follow.GetUserFansReq
	GetUserFansResp      = user_follow.GetUserFansResp
	GetUserFollowsReq    = user_follow.GetUserFollowsReq
	GetUserFollowsResp   = user_follow.GetUserFollowsResp
	IsUserFollowedReq    = user_follow.IsUserFollowedReq
	IsUserFollowedResp   = user_follow.IsUserFollowedResp
	UserInfo             = user_follow.UserInfo

	UserFollowService interface {
		// 创建用户关注关系
		CreateUserFollow(ctx context.Context, in *CreateUserFollowReq, opts ...grpc.CallOption) (*CreateUserFollowResp, error)
		// 删除用户关注关系
		DeleteUserFollow(ctx context.Context, in *DeleteUserFollowReq, opts ...grpc.CallOption) (*DeleteUserFollowResp, error)
		// 获取用户关注列表
		GetUserFollows(ctx context.Context, in *GetUserFollowsReq, opts ...grpc.CallOption) (*GetUserFollowsResp, error)
		// 获取用户粉丝列表
		GetUserFans(ctx context.Context, in *GetUserFansReq, opts ...grpc.CallOption) (*GetUserFansResp, error)
		// 检查是否关注
		IsUserFollowed(ctx context.Context, in *IsUserFollowedReq, opts ...grpc.CallOption) (*IsUserFollowedResp, error)
		// 获取互相关注列表
		GetMutualFollows(ctx context.Context, in *GetMutualFollowsReq, opts ...grpc.CallOption) (*GetMutualFollowsResp, error)
	}

	defaultUserFollowService struct {
		cli zrpc.Client
	}
)

func NewUserFollowService(cli zrpc.Client) UserFollowService {
	return &defaultUserFollowService{
		cli: cli,
	}
}

// 创建用户关注关系
func (m *defaultUserFollowService) CreateUserFollow(ctx context.Context, in *CreateUserFollowReq, opts ...grpc.CallOption) (*CreateUserFollowResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.CreateUserFollow(ctx, in, opts...)
}

// 删除用户关注关系
func (m *defaultUserFollowService) DeleteUserFollow(ctx context.Context, in *DeleteUserFollowReq, opts ...grpc.CallOption) (*DeleteUserFollowResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.DeleteUserFollow(ctx, in, opts...)
}

// 获取用户关注列表
func (m *defaultUserFollowService) GetUserFollows(ctx context.Context, in *GetUserFollowsReq, opts ...grpc.CallOption) (*GetUserFollowsResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.GetUserFollows(ctx, in, opts...)
}

// 获取用户粉丝列表
func (m *defaultUserFollowService) GetUserFans(ctx context.Context, in *GetUserFansReq, opts ...grpc.CallOption) (*GetUserFansResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.GetUserFans(ctx, in, opts...)
}

// 检查是否关注
func (m *defaultUserFollowService) IsUserFollowed(ctx context.Context, in *IsUserFollowedReq, opts ...grpc.CallOption) (*IsUserFollowedResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.IsUserFollowed(ctx, in, opts...)
}

// 获取互相关注列表
func (m *defaultUserFollowService) GetMutualFollows(ctx context.Context, in *GetMutualFollowsReq, opts ...grpc.CallOption) (*GetMutualFollowsResp, error) {
	client := user_follow.NewUserFollowServiceClient(m.cli.Conn())
	return client.GetMutualFollows(ctx, in, opts...)
}