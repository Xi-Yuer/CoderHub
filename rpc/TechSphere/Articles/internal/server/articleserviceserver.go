// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: articles.proto

package server

import (
	"context"

	"coderhub/rpc/TechSphere/Articles/articles"
	"coderhub/rpc/TechSphere/Articles/internal/logic"
	"coderhub/rpc/TechSphere/Articles/internal/svc"
)

type ArticleServiceServer struct {
	svcCtx *svc.ServiceContext
	articles.UnimplementedArticleServiceServer
}

func NewArticleServiceServer(svcCtx *svc.ServiceContext) *ArticleServiceServer {
	return &ArticleServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *ArticleServiceServer) GetArticle(ctx context.Context, in *articles.GetArticleRequest) (*articles.GetArticleResponse, error) {
	l := logic.NewGetArticleLogic(ctx, s.svcCtx)
	return l.GetArticle(in)
}

func (s *ArticleServiceServer) CreateArticle(ctx context.Context, in *articles.CreateArticleRequest) (*articles.CreateArticleResponse, error) {
	l := logic.NewCreateArticleLogic(ctx, s.svcCtx)
	return l.CreateArticle(in)
}

func (s *ArticleServiceServer) UpdateArticle(ctx context.Context, in *articles.UpdateArticleRequest) (*articles.UpdateArticleResponse, error) {
	l := logic.NewUpdateArticleLogic(ctx, s.svcCtx)
	return l.UpdateArticle(in)
}

func (s *ArticleServiceServer) UpdateLikeCount(ctx context.Context, in *articles.UpdateLikeCountRequest) (*articles.UpdateLikeCountResponse, error) {
	l := logic.NewUpdateLikeCountLogic(ctx, s.svcCtx)
	return l.UpdateLikeCount(in)
}

func (s *ArticleServiceServer) DeleteArticle(ctx context.Context, in *articles.DeleteArticleRequest) (*articles.DeleteArticleResponse, error) {
	l := logic.NewDeleteArticleLogic(ctx, s.svcCtx)
	return l.DeleteArticle(in)
}
