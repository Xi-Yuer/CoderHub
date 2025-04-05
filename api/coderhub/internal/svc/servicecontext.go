package svc

import (
	"coderhub/api/coderhub/internal/config"
	"coderhub/pkg/ws"
	"coderhub/repository"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/storage"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                      config.Config
	UserService                 coderhub.UserServiceClient
	ImageAuthService            coderhub.ImageServiceClient
	ArticlesService             coderhub.ArticleServiceClient
	AcademicService             coderhub.AcademicNavigatorServiceClient
	UserFollowService           coderhub.UserFollowServiceClient
	ImagesService               coderhub.ImageServiceClient
	CommentService              coderhub.CommentServiceClient
	QuestionBankService         coderhub.QuestionServiceClient
	FavoriteService             coderhub.FavorFoldServiceClient
	FavoriteContentService      coderhub.FavorServiceClient
	EmotionService              coderhub.EmotionServiceClient
	TagService                  coderhub.ArticleTagServiceClient
	QuestionBankCategoryService coderhub.QuestionBankCategoryServiceClient
	SchoolExpService            coderhub.SchoolExpServiceClient
	WorkExpService              coderhub.WorkExpServiceClient
	MessageService              coderhub.MessageServiceClient
	CreatorDashBoardService     coderhub.CreatorDashBoardServiceClient
	UserSessionRepository       repository.UserSessionRepository
	PrivateMessageRepository    repository.PrivateMessageRepository
	UserRepository              repository.UserRepository
	WsHub                       *ws.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	hub, err := ws.NewHub()
	if err != nil {
		panic(err)
	}
	// 启动 Hub 的核心管理逻辑
	go hub.Run()
	// 启动 WebSocket 连接的心跳检测
	go hub.StartHeartbeat(30*time.Second, 60*time.Second)
	sql := storage.NewGorm()
	rdb, err := storage.NewRedisDB(storage.DefaultConfig())
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:                      c,
		UserService:                 coderhub.NewUserServiceClient(zrpc.MustNewClient(c.UserService).Conn()),
		ImageAuthService:            coderhub.NewImageServiceClient(zrpc.MustNewClient(c.ImageAuthService).Conn()),
		ArticlesService:             coderhub.NewArticleServiceClient(zrpc.MustNewClient(c.ArticlesService).Conn()),
		AcademicService:             coderhub.NewAcademicNavigatorServiceClient(zrpc.MustNewClient(c.AcademicService).Conn()),
		UserFollowService:           coderhub.NewUserFollowServiceClient(zrpc.MustNewClient(c.UserFollowService).Conn()),
		ImagesService:               coderhub.NewImageServiceClient(zrpc.MustNewClient(c.ImagesService).Conn()),
		CommentService:              coderhub.NewCommentServiceClient(zrpc.MustNewClient(c.CommentService).Conn()),
		QuestionBankService:         coderhub.NewQuestionServiceClient(zrpc.MustNewClient(c.QuestionBankService).Conn()),
		FavoriteService:             coderhub.NewFavorFoldServiceClient(zrpc.MustNewClient(c.FavoriteService).Conn()),
		FavoriteContentService:      coderhub.NewFavorServiceClient(zrpc.MustNewClient(c.FavoriteContentService).Conn()),
		EmotionService:              coderhub.NewEmotionServiceClient(zrpc.MustNewClient(c.EmotionService).Conn()),
		TagService:                  coderhub.NewArticleTagServiceClient(zrpc.MustNewClient(c.TagService).Conn()),
		QuestionBankCategoryService: coderhub.NewQuestionBankCategoryServiceClient(zrpc.MustNewClient(c.QuestionBankCategoryService).Conn()),
		SchoolExpService:            coderhub.NewSchoolExpServiceClient(zrpc.MustNewClient(c.SchoolExpService).Conn()),
		WorkExpService:              coderhub.NewWorkExpServiceClient(zrpc.MustNewClient(c.WorkExpService).Conn()),
		MessageService:              coderhub.NewMessageServiceClient(zrpc.MustNewClient(c.MessageService).Conn()),
		CreatorDashBoardService:     coderhub.NewCreatorDashBoardServiceClient(zrpc.MustNewClient(c.CreatorDashBoardService).Conn()),
		PrivateMessageRepository:    repository.NewPrivateMessageRepository(sql, rdb),
		UserSessionRepository:       repository.NewUserSessionRepository(sql),
		UserRepository:              repository.NewUserRepositoryImpl(sql, rdb),
		WsHub:                       hub,
	}
}
