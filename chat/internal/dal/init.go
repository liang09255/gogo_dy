package dal

import (
	"chat/internal/database"
	"chat/internal/model"
	"common/ggLog"
)

func Init() {
	db := database.GetDB()
	err := db.AutoMigrate(
		&model.Message{},
	)
	if err != nil {
		ggLog.Fatalf("auto migrate error: %v", err)
	}
}
