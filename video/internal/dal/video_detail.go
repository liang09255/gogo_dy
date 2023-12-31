package dal

import (
	"common/ggLog"
	"context"
	"errors"
	"gorm.io/gorm"
	"video/internal/database"
	"video/internal/model"
	"video/internal/repo"
)

type VideoDetailDal struct {
	conn *database.GormConn
}

var _ repo.VideoDetailRepo = (*VideoDetailDal)(nil)

func NewVideoDetailDao() *VideoDetailDal {
	return &VideoDetailDal{
		conn: database.New(),
	}
}

// 插入记录,发布视频时就可以考虑插入一条记录
func (vd *VideoDetailDal) Insert(ctx context.Context, vid int64) error {
	v := &model.VideoDetail{
		Id:            vid,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Create(&v)
	return t.Error
}

// 获得视频的信息
func (vd *VideoDetailDal) GetVideoDetail(ctx context.Context, vid int64) (*model.VideoDetail, error) {
	v := &model.VideoDetail{}
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Where("id = ?", vid).
		First(&v)
	if t.Error != nil {
		return nil, t.Error
	}
	return v, nil
}

func (vd *VideoDetailDal) AddFavoriteCount(ctx context.Context, vid int64, count int64) error {
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Where("id = ?", vid).
		Update("favorite_count", gorm.Expr("favorite_count + ?", count))
	return t.Error
}

func (vd *VideoDetailDal) SubFavoriteCount(ctx context.Context, vid int64, count int64) error {
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Where("id = ?", vid).
		Update("favorite_count", gorm.Expr("favorite_count - ?", count))
	return t.Error
}

func (vd *VideoDetailDal) AddCommentCount(ctx context.Context, vid int64, count int64) error {
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Where("id = ?", vid).
		Update("comment_count", gorm.Expr("comment_count + ?", count))
	return t.Error
}

func (vd *VideoDetailDal) SubCommentCount(ctx context.Context, vid int64, count int64) error {
	t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
		Where("id = ?", vid).
		Update("comment_count", gorm.Expr("comment_count - ?", count))
	return t.Error
}

func (vd *VideoDetailDal) BatchInsertFavorite(ctx context.Context, m map[int64]int) (faultCount int, err error) {
	for vid, num := range m {
		condition := &model.VideoDetail{Id: vid}
		t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
			Where(condition).Update("favorite_count", gorm.Expr("favorite_count + ?", num))
		if t.RowsAffected == 0 || errors.Is(t.Error, gorm.ErrRecordNotFound) {
			ggLog.Debugf("vid:%d 不存在该记录", vid)
			// 说明不存在该记录，需要插入一条
			t = vd.conn.WithContext(ctx).Create(
				&model.VideoDetail{
					Id:            vid,
					FavoriteCount: int64(num),
				})
		}
		if t.Error != nil {
			faultCount++
			ggLog.Error("批量插入数据错误", t.Error, "错误值为:vid %d: num:%d", vid, num)
		}
	}
	return
}

func (vd *VideoDetailDal) BatchInsertComment(ctx context.Context, m map[int64]int) (faultCount int, err error) {
	if len(m) == 0 {
		return 0, nil
	}
	for vid, num := range m {
		condition := &model.VideoDetail{Id: vid}
		t := vd.conn.WithContext(ctx).Model(&model.VideoDetail{}).
			Where(condition).Update("comment_count", gorm.Expr("comment_count + ?", num))
		if t.RowsAffected == 0 || errors.Is(t.Error, gorm.ErrRecordNotFound) {
			ggLog.Debugf("vid:%d 不存在该记录", vid)
			// 说明不存在该记录，需要插入一条
			t = vd.conn.WithContext(ctx).Create(
				&model.VideoDetail{
					Id:           vid,
					CommentCount: int64(num),
				})
		}
		if t.Error != nil {
			faultCount++
			ggLog.Error("批量插入数据错误", t.Error, "错误值为:vid %d: num:%d", vid, num)
		}
	}
	return
}
