package dal

import (
	"common/ggLog"
	"video/internal/database"
	"video/internal/model"
)

func Init() {
	db := database.GetDB()
	err := db.AutoMigrate(
		&model.Comment{},
		&model.Video{},
		&model.Favorite{},
		&model.VideoDetail{},
	)
	if err != nil {
		ggLog.Fatalf("auto migrate error: %v", err)
	}
}
