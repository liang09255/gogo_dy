package dal

import (
	"gorm.io/gorm"
	"main/global"
	"time"
)

type videoDal struct{}

var VideoDal = &videoDal{}

// Video 视频
type Video struct {
	gorm.Model
	Id            int64  `gorm:"primaryKey"`
	AuthorId      int64  `gorm:"not null"`
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64
	CommentCount  int64
	Title         string `gorm:"not null"`
}

// AddVideo 上传视频
func (h *videoDal) AddVideo(userID int64, playUrl, coverUrl, title string) error {
	var video = Video{
		AuthorId:      userID,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	return global.MysqlDB.Create(&video).Error
}

// GetPublishList 获取某个用户发表过的视频列表
func (h *videoDal) GetPublishList(id int64) ([]Video, error) {
	var videos []Video
	t := global.MysqlDB.Where("author_id = ?", id).Find(&videos)
	return videos, t.Error
}

// GetFeedList 获取视频feed列表
func (h *videoDal) GetFeedList(latest int64) ([]Video, error) {
	// 传入的毫秒级别时间戳，需要进一步处理
	// 这里拆分是为了避免毫秒直接转秒级会丢失进度，所以分为秒和毫秒两部分分别处理
	seconds := latest / 1000
	milliseconds := latest % 1000
	var videos []Video
	t := global.MysqlDB.
		Where("created_at < ?", time.Unix(seconds, milliseconds*int64(time.Millisecond))).
		Order("created_at DESC").
		Limit(30).
		Find(&videos)
	return videos, t.Error
}

func (h *videoDal) MGetVideoInfo(ids []int64) ([]Video, error) {
	var videos []Video
	t := global.MysqlDB.Where("id in ?", ids).Find(&videos)
	return videos, t.Error
}

func (h *videoDal) AddFavoriteCount(videoID int64) error {
	return global.MysqlDB.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (h *videoDal) AddCommentCount(videoID int64) error {
	return global.MysqlDB.Model(&Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

func (h *videoDal) ReduceFavoriteCount(videoID int64) error {
	return global.MysqlDB.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

func (h *videoDal) ReduceCommentCount(videoID int64) error {
	return global.MysqlDB.Model(&Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}

func (h *videoDal) SelectByVId(videoID int64) (v *Video, err error) {
	err = global.MysqlDB.Model(&Video{}).Where("id = ?", videoID).Find(&v).Error
	return v, err
}
