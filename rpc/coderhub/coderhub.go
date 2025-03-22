package main

import (
	"flag"
	"fmt"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/config"
	academicnavigatorserviceServer "coderhub/rpc/coderhub/internal/server/academicnavigatorservice"
	articleserviceServer "coderhub/rpc/coderhub/internal/server/articleservice"
	articletagserviceServer "coderhub/rpc/coderhub/internal/server/articletagservice"
	commentserviceServer "coderhub/rpc/coderhub/internal/server/commentservice"
	creatorDashBoardService "coderhub/rpc/coderhub/internal/server/creatordashboardservice"
	emotionserviceServer "coderhub/rpc/coderhub/internal/server/emotionservice"
	favorfoldserviceServer "coderhub/rpc/coderhub/internal/server/favorfoldservice"
	favorserviceServer "coderhub/rpc/coderhub/internal/server/favorservice"
	imagerelationserviceServer "coderhub/rpc/coderhub/internal/server/imagerelationservice"
	imageserviceServer "coderhub/rpc/coderhub/internal/server/imageservice"
	messageserviceServer "coderhub/rpc/coderhub/internal/server/messageservice"
	questionbankcategoryserviceServe "coderhub/rpc/coderhub/internal/server/questionbankcategoryservice"
	questionserviceServer "coderhub/rpc/coderhub/internal/server/questionservice"
	schoolExpServer "coderhub/rpc/coderhub/internal/server/schoolexpservice"
	userfollowserviceServer "coderhub/rpc/coderhub/internal/server/userfollowservice"
	userserviceServer "coderhub/rpc/coderhub/internal/server/userservice"
	workExpServer "coderhub/rpc/coderhub/internal/server/workexpservice"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/coderhub.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		coderhub.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		coderhub.RegisterUserFollowServiceServer(grpcServer, userfollowserviceServer.NewUserFollowServiceServer(ctx))
		coderhub.RegisterAcademicNavigatorServiceServer(grpcServer, academicnavigatorserviceServer.NewAcademicNavigatorServiceServer(ctx))
		coderhub.RegisterArticleServiceServer(grpcServer, articleserviceServer.NewArticleServiceServer(ctx))
		coderhub.RegisterCommentServiceServer(grpcServer, commentserviceServer.NewCommentServiceServer(ctx))
		coderhub.RegisterImageRelationServiceServer(grpcServer, imagerelationserviceServer.NewImageRelationServiceServer(ctx))
		coderhub.RegisterImageServiceServer(grpcServer, imageserviceServer.NewImageServiceServer(ctx))
		coderhub.RegisterQuestionServiceServer(grpcServer, questionserviceServer.NewQuestionServiceServer(ctx))
		coderhub.RegisterFavorFoldServiceServer(grpcServer, favorfoldserviceServer.NewFavorFoldServiceServer(ctx))
		coderhub.RegisterFavorServiceServer(grpcServer, favorserviceServer.NewFavorServiceServer(ctx))
		coderhub.RegisterEmotionServiceServer(grpcServer, emotionserviceServer.NewEmotionServiceServer(ctx))
		coderhub.RegisterArticleTagServiceServer(grpcServer, articletagserviceServer.NewArticleTagServiceServer(ctx))
		coderhub.RegisterQuestionBankCategoryServiceServer(grpcServer, questionbankcategoryserviceServe.NewQuestionBankCategoryServiceServer(ctx))
		coderhub.RegisterSchoolExpServiceServer(grpcServer, schoolExpServer.NewSchoolExpServiceServer(ctx))
		coderhub.RegisterWorkExpServiceServer(grpcServer, workExpServer.NewWorkExpServiceServer(ctx))
		coderhub.RegisterMessageServiceServer(grpcServer, messageserviceServer.NewMessageServiceServer(ctx))
		coderhub.RegisterCreatorDashBoardServiceServer(grpcServer, creatorDashBoardService.NewCreatorDashBoardServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
