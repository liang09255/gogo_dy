package baseCtlModel

type APIError struct {
	StatusCode int32
	StatusMsg  string
}
type APIBaseResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func newError(code int32, msg string) APIError {
	return APIError{
		StatusCode: code,
		StatusMsg:  msg,
	}
}

func newBaseResp(code int32, msg string) APIBaseResp {
	return APIBaseResp{
		StatusCode: code,
		StatusMsg:  msg,
	}
}

func NewBaseSuccessResp(msg ...string) APIBaseResp {
	if len(msg) == 0 {
		return Success
	}
	return newBaseResp(0, msg[0])
}
func (b APIError) Error() string {
	return b.StatusMsg
}
func (b APIError) WithDetails(msg string) APIError {
	return APIError{
		StatusCode: b.StatusCode,
		StatusMsg:  b.StatusMsg + ": " + msg,
	}
}
