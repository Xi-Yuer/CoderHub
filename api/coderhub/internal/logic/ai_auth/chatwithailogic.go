package ai_auth

import (
	"bufio"
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

// ChatWithAI 流式返回 AI 对话
func (l *ChatWithAILogic) ChatWithAI(w http.ResponseWriter, req *types.ChatWithAIReq) {
	url := "https://api.deepseek.com/chat/completions"
	request, _ := http.NewRequest("POST", url, strings.NewReader(utils.MessageHandler(req.Content)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Transfer-Encoding", "chunked")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.DeepSeekSearchKey))

	// 获取响应
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, "AI 服务器请求失败", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "服务器不支持流式响应", http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(res.Body)
	scanner := bufio.NewScanner(reader)
	isFirstChunk := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		// 解析 AI API 响应
		var response ChatResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			if content != "" {
				status := 1 // 默认中间部分
				if isFirstChunk {
					status = 0 // 第一部分
					isFirstChunk = false
				}

				chatResponse := ChatStreamResponse{
					Code:    200,
					Message: "success",
					Data: struct {
						Content string `json:"content"`
						Status  int    `json:"status"`
					}{content, status},
				}

				jsonData, _ := json.Marshal(chatResponse)

				// 关键修改：确保符合 SSE 规范
				_, _ = fmt.Fprintf(w, "data: %s\n\n", jsonData)
				flusher.Flush()
			}
		}
	}

	// 结束事件
	finalResponse := ChatStreamResponse{
		Code:    200,
		Message: "completed",
		Data: struct {
			Content string `json:"content"`
			Status  int    `json:"status"`
		}{"", 2}, // 2 代表流结束
	}
	jsonData, _ := json.Marshal(finalResponse)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}

type ChatStreamResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Content string `json:"content"`
		Status  int    `json:"status"`
	} `json:"data"`
}

type ChatCompletionChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type ChatResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}
