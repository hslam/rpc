package cons

const (
	OK					= 10200
	BadRequest			= 10400
	Unauthorized		= 10401
	InvalidCredential	= 10402
	PermissionDenied	= 10403
	NotFound 			= 10404
	InternalError		= 10500
	ServiceUnavailable	= 10503
	TimedOut			= 10504
	TooManyRequests		= 10505
	Conflict			= 10506
	UnknownError		= 10600

	NetworkFailure		= 10601
	CookieExpired		= 10700
	RuntimeError		= 10800

	// cache specific
	CacheGetFailed	= 20404
	CacheSetFailed		= 20405

	// Server specific
)

var statusText = map[int]string{
	OK:                 "成功",
	BadRequest:         "请求错误",
	Unauthorized:       "未授权的访问",
	InvalidCredential:  "无效的授权",
	PermissionDenied:   "权限错误",
	NotFound:           "未找到",
	InternalError:      "服务内部错误",
	ServiceUnavailable: "服务不可用",
	TimedOut:           "处理超时",
	TooManyRequests:    "请求过于频繁",
	Conflict:           "资源冲突",
	UnknownError:       "未知错误",

	NetworkFailure:		"网络错误",
	CookieExpired:		"Cookie失效",
	RuntimeError:		"运行时错误",

	// for cache
	CacheGetFailed:		"缓存读取失败",
	CacheSetFailed:		"缓存写入失败",

}

// Msg returns a text for the the status code. It returns the empty
// string if the code is unknown.
func Msg(code int) string {
	return statusText[code]
}
