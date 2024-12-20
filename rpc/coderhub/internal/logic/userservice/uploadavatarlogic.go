package userservicelogic

import (
	"coderhub/model"
	imagerelationservicelogic "coderhub/rpc/coderhub/internal/logic/imagerelationservice"
	imageservicelogic "coderhub/rpc/coderhub/internal/logic/imageservice"
	"coderhub/shared/utils"
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UploadAvatar 上传用户头像
func (l *UploadAvatarLogic) UploadAvatar(in *coderhub.UploadAvatarRequest) (*coderhub.UploadAvatarResponse, error) {
	// 权限校验
	var (
		userId string
		err    error
	)
	// 从 metadata 中获取 userId
	if userId, err = utils.GetUserMetaData(l.ctx); err != nil {
		return nil, err
	}

	if userId != strconv.FormatInt(in.UserId, 10) {
		return nil, fmt.Errorf("非法操作")
	}

	// 保存图片关系
	createRelationRequest := &coderhub.CreateRelationRequest{
		ImageId:    in.ImageId,
		EntityId:   in.UserId,
		EntityType: model.ImageRelationUserAvatar,
	}
	batchCreateRelationService := imagerelationservicelogic.NewBatchCreateRelationLogic(l.ctx, l.svcCtx)
	_, err = batchCreateRelationService.BatchCreateRelation(
		&coderhub.BatchCreateRelationRequest{
			Relations: []*coderhub.CreateRelationRequest{
				createRelationRequest,
			},
		},
	)
	if err != nil {
		// 事务回滚
		batchDeleteRelationService := imagerelationservicelogic.NewBatchDeleteRelationLogic(l.ctx, l.svcCtx)
		_, err := batchDeleteRelationService.BatchDeleteRelation(&coderhub.BatchDeleteRelationRequest{
			Ids: []int64{in.ImageId},
		})
		if err != nil {
			return nil, err
		}
		l.Errorf("保存图片关系失败: %v", err)
		return nil, err
	}
	// 获取图片信息
	imageInfoService := imageservicelogic.NewGetLogic(l.ctx, l.svcCtx)
	imageInfo, err := imageInfoService.Get(&coderhub.GetRequest{
		ImageId: in.ImageId,
	})
	if err != nil {
		return nil, err
	}
	// 更新用户头像
	userInfo, err := l.svcCtx.UserRepository.GetUserByID(in.UserId)
	if err != nil {
		return nil, err
	}
	l.Logger.Infof("imageRelation用户头像: %v", imageInfo.Url)
	userInfo.Avatar = sql.NullString{String: imageInfo.Url, Valid: true}
	l.Logger.Infof("更新用户头像: %v", userInfo.Avatar)
	err = l.svcCtx.UserRepository.UpdateUser(userInfo)
	if err != nil {
		l.Errorf("更新用户头像失败: %v", err)
		return nil, err
	}

	return &coderhub.UploadAvatarResponse{
		ImageId:      in.ImageId,
		Url:          imageInfo.Url,
		ThumbnailUrl: imageInfo.ThumbnailUrl,
	}, nil
}
