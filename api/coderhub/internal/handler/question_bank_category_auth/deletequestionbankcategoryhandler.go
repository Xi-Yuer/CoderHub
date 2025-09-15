package question_bank_category_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/question_bank_category_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除题库分类
func DeleteQuestionBankCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteQuestionBankCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question_bank_category_auth.NewDeleteQuestionBankCategoryLogic(r.Context(), svcCtx)
		resp, err := l.DeleteQuestionBankCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
