package dal

import (
	"main/global"
	"main/model"
)

type userDal struct{}

var UserDal = &userDal{}

func (h *userDal) Exist(username string) error {
	var db = global.MysqlDB
	var user = model.User{Username: username}
	t := db.First(&user)
	return t.Error
}

func (h *userDal) Adduser(username string, password string) error {
	var db = global.MysqlDB
	var user = model.User{
		Username:      username,
		Password:      password,
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
	}
	t := db.Create(&user)
	return t.Error
}
