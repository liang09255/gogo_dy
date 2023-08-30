package dal

import (
	"common/ggLog"
	"user/internal/database"
	"user/internal/model"
)

func Init() {
	db := database.GetDB()
	err := db.AutoMigrate(
		&model.User{},
		&model.Relation{},
	)
	if err != nil {
		ggLog.Fatalf("auto migrate error: %v", err)
	}
}
