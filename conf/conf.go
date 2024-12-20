package conf

var JWTSecret = "secret_key"

var GoMail = struct {
	Username string
	Password string
	Host     string
	Port     int
}{
	Username: "2214380963@qq.com",
	Password: "uylmxoxamqsddjfc",
	Host:     "smtp.qq.com",
	Port:     587,
}

var HttpCode = struct {
	HttpStatusOK      int32
	HttpForbidden     int32
	HttpBadRequest    int32
	HttpNotFound      int32
	HttpParamsMissing int32
}{
	HttpStatusOK:      200,
	HttpForbidden:     403,
	HttpBadRequest:    400,
	HttpNotFound:      404,
	HttpParamsMissing: 422,
}
var HttpMessage = struct {
	MsgOK     string
	MsgError  string
	MsgFailed string
}{
	MsgOK:     "OK",
	MsgError:  "Error",
	MsgFailed: "Failed",
}

const (
	MaxImageCount = 20
	UrlHTTP       = "http://"
	UrlHTTPS      = "https://"
)

// 错误常量
const (
	ErrValidationFailed    = "参数校验失败: %v"
	ErrUserMetaFailed      = "获取用户元数据失败: %v"
	ErrUserIDConversion    = "转换用户ID失败: %v"
	ErrArticleNotFound     = "文章不存在"
	ErrNoPermission        = "您无权修改此文章"
	ErrUpdateFailed        = "更新文章失败: %v"
	ErrImageCountExceeded  = "图片数量不能超过%d张"
	ErrInvalidImageURL     = "图片URL格式不正确"
	ErrInvalidCoverURL     = "封面图URL格式不正确"
	ErrGetArticleFailed    = "获取文章失败: %v"
	ErrDeleteArticleFailed = "删除文章失败: %v"
)
