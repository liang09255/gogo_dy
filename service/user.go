package service

import "main/dal"

type registerinfo struct{}

var Registerinfo *registerinfo

func (h *registerinfo) Register(username string, password string) error {

	Dal := dal.UserDal
	//判断用户名是否已经存在
	err := Dal.Exist(username)

	//若存在,则直接返回
	if err != nil {
		return err
	}
	//若不存在，则将注册信息写入数据库中
	err = Dal.Adduser(username, password)

	if err != nil {
		return err
	}

	return nil
}

func (h *registerinfo) Login(username string, password string) (msg string, err error) {

	return
}
