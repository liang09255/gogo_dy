package dal

import (
	"video/internal/database"
	"video/internal/repo"
)

type TranDal struct{}

var _ repo.TranRepo = (*TranDal)(nil)

func NewTranRepo() *TranDal {
	return &TranDal{}
}

func (t TranDal) NewTransactionConn() database.DbConn {
	return database.New()
}
