package common

type ResponseMsg struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ErrorMsg ResponseMsg

const (
	ParamErrorCode         int = 60400
	LoginErrorCode         int = 60401
	TokenGenerateErrorCode int = 60402
	TokenInvalidErrorCode  int = 60403
	ServerErrorCode        int = 60500

	SuccessCode int = 0
)

type SuccessMsg ResponseMsg

var (
	LoginError = ErrorMsg{
		Code:    LoginErrorCode,
		Message: "用户名和密码是无效的",
	}

	TokenGenerateError = ErrorMsg{
		Code:    TokenGenerateErrorCode,
		Message: "生成令牌失败",
	}
	TokenInvalidError = ErrorMsg{
		Code:    TokenInvalidErrorCode,
		Message: "令牌无效",
	}

	TokenInfoGetError = ErrorMsg{
		Code:    TokenInvalidErrorCode,
		Message: "令牌信息获取失败",
	}

	ParamError = ErrorMsg{
		Code:    ParamErrorCode,
		Message: "参数错误，检查后重试",
	}
	NoUpdateError = ErrorMsg{
		Code:    ParamErrorCode,
		Message: "参数没有任何改变，检查后重试",
	}
	ExistsError = ErrorMsg{
		Code:    ParamErrorCode,
		Message: "对象已经存在，检查后重试",
	}
	NotExistsError = ErrorMsg{
		Code:    ParamErrorCode,
		Message: "对象不存在，检查后重试",
	}

	ServerError = ErrorMsg{
		Code:    ServerErrorCode,
		Message: "服务器内部错误",
	}

	CreateError = ErrorMsg{
		Code:    ServerErrorCode,
		Message: "对象添加失败，等待后重试",
	}

	UpdateError = ErrorMsg{
		Code:    ServerErrorCode,
		Message: "对象更新失败，等待后重试",
	}

	DeleteError = ErrorMsg{
		Code:    ServerErrorCode,
		Message: "对象删除失败，等待后重试",
	}
	NotEmptyError = ErrorMsg{
		Code:    ParamErrorCode,
		Message: "存在关联的文章，无法删除",
	}
)
