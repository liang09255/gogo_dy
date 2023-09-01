package dal

import (
	"context"
	"time"
	"video/internal/database"
	"video/internal/model"
	"video/internal/repo"
)

type VideoDal struct {
	conn *database.GormConn
}

var _ repo.VideoRepo = (*VideoDal)(nil)

func NewVideoDao() *VideoDal {
	return &VideoDal{
		conn: database.New(),
	}
}

func (v *VideoDal) NewVideo(ctx context.Context, userID int64, videoURL, coverURL, title string) error {
	var video = model.Video{
		AuthorId:      userID,
		PlayUrl:       videoURL,
		CoverUrl:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	return v.conn.WithContext(ctx).Create(&video).Error
}

func (v *VideoDal) PublishList(ctx context.Context, userID int64) ([]*model.Video, error) {
	var videos []*model.Video
	t := v.conn.WithContext(ctx).Where("author_id = ?", userID).Find(&videos)
	return videos, t.Error
}

func (v *VideoDal) Feed(ctx context.Context, latest int64) ([]*model.Video, error) {
	var videos []*model.Video
	// 需要查找比latest早发布的视频
	if latest == 0 {
		latest = time.Now().UnixMicro()
	}
	t := v.conn.WithContext(ctx).
		Where("created_at < ?", time.UnixMicro(latest)).
		Order("created_at DESC").
		Limit(30).
		Find(&videos)
	return videos, t.Error
}

// MGetVideoInfo 批量获得视频信息
func (v *VideoDal) MGetVideoInfo(ctx context.Context, ids []int64) ([]*model.Video, error) {
	var videoList []*model.Video
	t := v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id in ?", ids).
		Find(&videoList)

	return videoList, t.Error
}

func (v *VideoDal) GetVideoInfo(ctx context.Context, videoId int64) (*model.Video, error) {
	var video *model.Video
	t := v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id = ?", videoId).
		Find(&video)

	return video, t.Error
}
