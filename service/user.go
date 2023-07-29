package service

type registerinfo struct{}

var Registerinfo *registerinfo

func (h *registerinfo) Register(username string, password string) (msg string, err error) {
	return
}
