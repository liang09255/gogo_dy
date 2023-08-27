package dal

import (
	"chat/internal/database"
	"chat/internal/repo"
)

type TranDal struct{}

var _ repo.TranRepo = (*TranDal)(nil)

func NewTranRepo() *TranDal {
	return &TranDal{}
}

func (t TranDal) NewTransactionConn() database.DbConn {
	return database.New()
}
