package repo

import "chat/internal/database"

type TranRepo interface {
	NewTransactionConn() database.DbConn
}
