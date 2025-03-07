package ai_auth

import (
	"coderhub/conf"
	"coderhub/shared/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWithAILogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChatWithAILogic 与AI对话
func NewChatWithAILogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWithAILogic {
	return &ChatWithAILogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatWithAILogic) ChatWithAI(req *types.ChatWithAIReq) (resp *types.ChatWithAIResp, err error) {
	url := "https://api.deepseek.com/chat/completions"
	method := "POST"

	payload := strings.NewReader(utils.MessageHandler(req.Content))

	client := &http.Client{}
	request, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.DeepSeekSearchKey))

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result ChatCompletion

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return &types.ChatWithAIResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: "",
		}, nil
	}
	return &types.ChatWithAIResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: result.Choices[0].Message.Content,
	}, nil
}

// ChatCompletion 定义 JSON 对应的结构体
type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     *string `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens          int                 `json:"prompt_tokens"`
	CompletionTokens      int                 `json:"completion_tokens"`
	TotalTokens           int                 `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int                 `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int                 `json:"prompt_cache_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}
