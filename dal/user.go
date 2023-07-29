package dal

import (
	"main/global"
	"main/model"
)

type userDal struct{}

var UserDal = &userDal{}

func (h *hellDal) Exist(username string) error {
	var db = global.MysqlDB
	var user = model.User{Username: username}
	t := db.First(&user)
	return t.Error
}
