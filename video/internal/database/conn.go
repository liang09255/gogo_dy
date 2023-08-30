package database

type DbConn interface {
	Begin()
	Rollback()
	Commit()
}
