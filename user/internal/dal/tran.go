package dal

import (
	"user/internal/database"
	"user/internal/repo"
)

type TranDal struct{}

var _ repo.TranRepo = (*TranDal)(nil)

func NewTranRepo() *TranDal {
	return &TranDal{}
}

func (t TranDal) NewTransactionConn() database.DbConn {
	return database.New()
}
