package tran

import "user/internal/database"

type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
