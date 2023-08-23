package baseCtlModel

const (
	SuccessCode = 0
	FailCode    = 1
)

type BaseResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func NewBaseSuccessResp(msg ...string) BaseResp {
	if len(msg) == 0 {
		return BaseResp{
			StatusCode: SuccessCode,
			StatusMsg:  "Success",
		}
	}
	return BaseResp{
		StatusCode: SuccessCode,
		StatusMsg:  msg[0],
	}
}

func NewBaseFailedResp(msg ...string) BaseResp {
	if len(msg) == 0 {
		return BaseResp{
			StatusCode: FailCode,
			StatusMsg:  "Failed",
		}
	}
	return BaseResp{
		StatusCode: FailCode,
		StatusMsg:  msg[0],
	}
}
