package service

import (
	"fmt"
	"main/dal"
)

type helloService struct{}

var HelloService = &helloService{}

func (h *helloService) SayHello(name string) (msg string, err error) {
	err = dal.HelloDal.NewHello(name)
	if err != nil {
		return
	}
	// FIXME 暂未实现获取用户id功能
	msg = fmt.Sprintf("hello %s, your id is %s", name, "1")
	return
}
