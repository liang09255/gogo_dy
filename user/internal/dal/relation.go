package dal

import (
	"context"
	"user/internal/database"
	"user/internal/model"
	"user/internal/repo"
)

type RelationDal struct {
	conn *database.GormConn
}

var _ repo.RelationRepo = (*RelationDal)(nil)

func NewRelationDao() *RelationDal {
	return &RelationDal{
		conn: database.New(),
	}
}

func (rd *RelationDal) Follow(ctx context.Context, uid int64, targetId int64) error {
	res := model.Relation{
		FollowId: uid,
		UserId:   targetId,
	}

	return rd.conn.WithContext(ctx).Create(&res).Error
}

func (rd *RelationDal) Delete(ctx context.Context, uid int64, targetId int64) error {
	return rd.conn.WithContext(ctx).Where("follow_id = ? AND user_id = ?", uid, targetId).Delete(&model.Relation{}).Error
}

// GetAllFollow 获取关注列表
func (rd *RelationDal) GetAllFollow(ctx context.Context, id int64) (followIDs []int64, err error) {
	var follows []model.Relation
	err = rd.conn.WithContext(ctx).Where("follow_id = ?", id).Find(&follows).Error

	if err != nil {
		return
	}
	for _, follow := range follows {
		followIDs = append(followIDs, follow.UserId)
	}
	return
}

// GetAllFollower 获取粉丝列表
func (rd *RelationDal) GetAllFollower(ctx context.Context, id int64) (followerIds []int64, err error) {
	var followers []model.Relation
	err = rd.conn.WithContext(ctx).Where("user_id = ?", id).Find(&followers).Limit(100).Error
	if err != nil {
		return
	}
	for _, follower := range followers {
		followerIds = append(followerIds, follower.FollowId)
	}
	return
}

func (rd *RelationDal) GetAllFriend(ctx context.Context, id int64) (friendIds []int64, err error) {
	var followIds []int64
	// 我关注的
	followIds, err = rd.GetAllFollow(ctx, id)
	if err != nil {
		return
	}
	var friends []model.Relation
	// 关注我的
	err = rd.conn.WithContext(ctx).Where("user_id = ? AND follow_id IN ?", id, followIds).Select("follow_id").Find(&friends).Error
	if err != nil {
		return
	}
	for _, friend := range friends {
		friendIds = append(friendIds, friend.FollowId)
	}
	return
}

func (rd *RelationDal) IsFollow(ctx context.Context, uid int64, targetId int64) (bool, error) {
	var count int64
	err := rd.conn.WithContext(ctx).Model(&model.Relation{}).Where("follow_id = ? AND user_id = ?", uid, targetId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
