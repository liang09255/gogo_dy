package repo

import "video/internal/database"

type TranRepo interface {
	NewTransactionConn() database.DbConn
}
