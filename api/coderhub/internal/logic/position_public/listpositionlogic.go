package position_public

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListPositionLogic 获取职位列表
func NewListPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPositionLogic {
	return &ListPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPositionLogic) ListPosition(req *types.GetPositionListReq) (resp *types.GetPositionListResp, err error) {
	// 发送 HTTP GET 请求
	httpResp, err := http.Get("https://aab0.github.io/data/new.json")
	if err != nil {
		return &types.GetPositionListResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: "http请求失败",
			},
			Data: types.GetPositionListRes{},
		}, nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResp.Body)

	// 检查 HTTP 响应状态码
	if httpResp.StatusCode != http.StatusOK {
		return &types.GetPositionListResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: "状态码有误",
			},
			Data: types.GetPositionListRes{},
		}, nil
	}

	// 定义一个临时变量用于解析 HTTP 响应
	var data types.GetPositionListRes
	if err := json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		return &types.GetPositionListResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: "解析 JSON 失败",
			},
			Data: types.GetPositionListRes{},
		}, nil
	}

	// 返回成功的响应
	return &types.GetPositionListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "成功",
		},
		Data: data, // 这里直接赋值解析后的数据
	}, nil
}
