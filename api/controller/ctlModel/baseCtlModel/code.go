package baseCtlModel

// 0 Success
var Success = newBaseResp(0, "Success")

// 400 BAD REQUEST
var (
	InvalidParams         = newError(40001, "参数错误")
	CommentActionUnknown  = newError(40003, "未知的评论操作")
	FavoriteActionUnknown = newError(40004, "未知的收藏操作")
	FileUploadFailed      = newError(40005, "文件上传失败")
	FileIsNotVideo        = newError(40006, "文件不是视频")
)

// 401 WITHOUT PERMISSION
var (
	InvalidToken = newError(40102, "无效的Token")
)

// 404 NOT FOUND
var (
	UserNotFound = newError(40401, "用户不存在")
)

// 500 INTERNAL ERROR
var (
	ServerInternal = newError(50001, "服务器内部错误")
)
