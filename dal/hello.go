package dal

import (
	"gorm.io/gorm"
	"main/global"
)

type Hello struct {
	gorm.Model
	Name string
}

type hellDal struct{}

var HelloDal = &hellDal{}

func (h *hellDal) NewHello(name string) error {
	t := global.MysqlDB.Create(&Hello{
		Name: name,
	})
	return t.Error
}
