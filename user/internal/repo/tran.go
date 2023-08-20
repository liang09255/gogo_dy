package repo

import "user/internal/database"

type TranRepo interface {
	NewTransactionConn() database.DbConn
}
