package svc

import (
	"coderhub/repository"
	"coderhub/rpc/coderhub/internal/config"
	"coderhub/shared/messaging"
	"coderhub/shared/storage"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type ServiceContext struct {
	Config                         config.Config
	Minio                          *storage.Minio
	UserRepository                 repository.UserRepository
	UserFollowRepository           repository.UserFollowRepository
	RedisDB                        storage.RedisDB
	GoMail                         messaging.GoMailImpl
	ImageRepository                repository.ImageRepository
	ImageRelationRepository        repository.ImageRelationRepository
	AcademicNavigatorRepository    repository.AcademicNavigatorRepository
	AcademicRelationLikeRepository repository.AcademicRelationLikeRepository
	ArticleRepository              repository.ArticleRepository
	ArticlesRelationLikeRepository repository.ArticlesRelationLikeRepository
	ArticlePVRepository            repository.ArticlePVRepository
	CommentRepository              repository.CommentRepository
	CommentRelationLikeRepository  repository.CommentRelationLikeRepository
	QuestionBankRepository         repository.QuestionBankRepository
	QuestionRepository             repository.QuestionRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	minioClient := storage.NewMinio(
		c.Minio.Endpoint,
		c.Minio.AccessKey,
		c.Minio.SecretKey,
		c.Minio.Bucket,
		c.Minio.Region,
		c.Minio.UseSSL,
	)

	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
		Username:  "elastic",
		Password:  "2214380963Wx!!",
	}

	err := minioClient.Connect()
	if err != nil {
		panic(fmt.Sprintf("Minio 连接失败: %v", err))
	}
	redisDB, err := storage.NewRedisDB(storage.DefaultConfig())
	if err != nil {
		panic(err)
	}

	sql := storage.NewGorm()

	return &ServiceContext{
		Config:                         c,
		Minio:                          minioClient,
		RedisDB:                        redisDB,
		GoMail:                         messaging.NewGoMail(),
		ImageRepository:                repository.NewImageRepository(sql),
		ImageRelationRepository:        repository.NewImageRelationRepository(sql),
		AcademicRelationLikeRepository: repository.NewAcademicRelationLikeRepositoryImpl(sql),
		ArticleRepository:              repository.NewArticleRepositoryImpl(sql, redisDB),
		ArticlesRelationLikeRepository: repository.NewArticlesRelationLikeRepository(sql, redisDB),
		ArticlePVRepository:            repository.NewArticlePVRepositoryImpl(sql, redisDB),
		CommentRepository:              repository.NewCommentRepository(sql, redisDB),
		CommentRelationLikeRepository:  repository.NewCommentRelationLikeRepository(sql, redisDB),
		UserRepository:                 repository.NewUserRepositoryImpl(sql, redisDB),
		UserFollowRepository:           repository.NewUserFollowRepositoryImpl(sql, redisDB),
		AcademicNavigatorRepository:    repository.NewAcademicNavigatorRepositoryImpl(sql, &cfg),
		QuestionBankRepository:         repository.NewQuestionRepositoryRepositoryImpl(sql, redisDB),
		QuestionRepository:             repository.NewQuestionRepositoryImpl(sql, redisDB),
	}
}
