package question_bank_category_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/question_bank_category_public"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取题库分类列表
func ListQuestionBankCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetQuestionBankCategoryListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question_bank_category_public.NewListQuestionBankCategoryLogic(r.Context(), svcCtx)
		resp, err := l.ListQuestionBankCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
