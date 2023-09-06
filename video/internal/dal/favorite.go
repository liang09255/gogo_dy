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

type FavoriteDal struct {
	conn *database.GormConn
}

var _ repo.FavoriteRepo = (*FavoriteDal)(nil)

func NewFavoriteDao() *FavoriteDal {
	return &FavoriteDal{
		conn: database.New(),
	}
}

// 获得点赞列表
func (f *FavoriteDal) GetFavoriteList(ctx context.Context, userId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := f.conn.WithContext(ctx).
		Model(&model.Favorite{}).
		Select("video_id").Where("user_id=? and deleted_at is null", userId).
		Find(&ids).Error
	return ids, err
}

// 点赞
func (f *FavoriteDal) PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	favorite := model.Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	var id = -1
	// 这里改为使用first，原先用count要查全表，这里只需要找到第一个匹配项
	err := f.conn.WithContext(ctx).Model(model.Favorite{}).
		Select("id").
		Where("user_id = ? and video_id = ? and deleted_at is null ", userId, videoId).
		First(&id).Error
	if err == nil && id != -1 {
		// 说明存在该记录
		// 重复记录需要返回一个错误
		return errors.New("重复记录")
		//return nil
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return f.conn.WithContext(ctx).Model(&model.Favorite{}).Create(&favorite).Error
	}
	if err != nil {
		return err
	}

	return nil
}

// 取消点赞
func (f *FavoriteDal) CancelFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	// 存在记录则删除，不存在则直接返回
	// user_id和video_id可以考虑建一个联合索引
	var id = -1
	err := f.conn.WithContext(ctx).Model(model.Favorite{}).
		Select("id").
		Where("user_id = ? and video_id = ? and deleted_at is null ", userId, videoId).
		First(&id).Error
	// 不存在则直接返回，需要报错,以使得另外一边得以回滚，以免点赞表未修改，而其余表改动了导致数据不一致
	//if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	//	return error
	//}

	if err != nil {
		return err
	}

	// 这里用的是gorm的delete，因为有delete_at字段，自动实现软删除,如果需要直接删除记录还需要自己修改,要么删除delete_at字段，要么使用unscoped
	err = f.conn.WithContext(ctx).Model(model.Favorite{}).
		Where("user_id = ? and video_id = ?", userId, videoId).
		Delete(&model.Favorite{}).Error

	return err
}

// GetFavoriteListByUserId 根据用户id获得点赞视频列表
func (f *FavoriteDal) GetFavoriteListByUserId(ctx context.Context, userid int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := f.conn.WithContext(ctx).Model(&model.Favorite{}).
		Select("video_id").Where("user_id = ? and deleted_at is null", userid).
		Find(&ids).Error

	return ids, err
}

// GetFavoritedCount 获得用户点赞数
func (f *FavoriteDal) GetFavoriteCount(ctx context.Context, userID int64) int64 {
	var count int64
	err := f.conn.WithContext(ctx).Model(&model.Favorite{}).
		Where("user_id = ? and deleted_at is null", userID).
		Count(&count).Error
	if err != nil {
		ggLog.Error("GetFavoriteCount error:", err)
	}
	return count
}

// GetFavoriteCountByVideoID 获取视频被点赞数
func (f *FavoriteDal) GetFavoriteCountByVideoID(ctx context.Context, videoID []int64) map[int64]int64 {
	res := make(map[int64]int64)
	queryRes := make([]struct {
		VideoId int64
		Count   int64
	}, 0)
	err := f.conn.WithContext(ctx).Model(&model.Favorite{}).
		Select("video_id, count(*) as count").
		Where("video_id in ?", videoID).
		Group("video_id").
		Scan(&queryRes).Error
	if err != nil {
		ggLog.Errorf("获取点赞数错误:%v", err)
		return nil
	}
	for _, v := range queryRes {
		res[v.VideoId] = v.Count
	}
	return res
}
