package sandpack_public

import (
	"coderhub/conf"
	"net/http"

	"coderhub/api/coderhub/internal/logic/sandpack_public"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ListSandpackHandler 获取sandpack列表

type File struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type Data struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Template    string           `json:"template"`
	UserID      string           `json:"user_id"`
	ArticleID   string           `json:"article_id"`
	Files       map[string]*File `json:"files"`
	CreatedAt   int64            `json:"created_at"`
	UpdatedAt   int64            `json:"updated_at"`
}

type Response struct {
	Code    int32       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ListSandpackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSandpackProjectListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sandpack_public.NewListSandpackLogic(r.Context(), svcCtx)
		resp, err := l.ListSandpack(&req)

		if resp.Data == nil {
			result := Response{
				Code:    conf.HttpCode.HttpStatusOK,
				Message: conf.HttpMessage.MsgOK,
				Data:    nil,
			}
			httpx.OkJsonCtx(r.Context(), w, result)
			return
		}

		files := make(map[string]*File)
		for _, file := range resp.Data.Files {
			files[file.Name] = &File{
				Code:     file.Code,
				Language: file.Language,
				Name:     file.Name,
			}
		}
		var data = Data{
			ArticleID:   resp.Data.ArticleID,
			CreatedAt:   resp.Data.CreatedAt,
			Description: resp.Data.Description,
			Files:       files,
			ID:          resp.Data.ID,
			Name:        resp.Data.Name,
			Template:    resp.Data.Template,
			UpdatedAt:   resp.Data.UpdatedAt,
			UserID:      resp.Data.UserID,
		}

		result := Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
			Data:    data,
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, result)
		}
	}
}
