package question_bank_category_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/question_bank_category_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建题库分类
func CreateQuestionBankCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateQuestionBankCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question_bank_category_auth.NewCreateQuestionBankCategoryLogic(r.Context(), svcCtx)
		resp, err := l.CreateQuestionBankCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
