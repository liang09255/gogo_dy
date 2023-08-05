package service

import (
	"main/dal"
	"main/model"
)

type userinfo struct{}

var Userinfo *userinfo

func (h *userinfo) Register(username string, password string) error {

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

func (h *userinfo) Login(username string, password string) (int64, error) {
	Dal := dal.UserDal

	Id, err := Dal.GetUserID(username, password)

	return Id, err
}

func (h *userinfo) GetUserInfo(token string, userid string, resp *model.Userinfo) error {
	err := dal.UserDal.GetUserInfo(token, userid, resp)

	if err != nil {
		return err
	}
	return err
}
