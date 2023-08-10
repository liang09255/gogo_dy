package dal

import (
	"gorm.io/gorm"
	"main/global"
)

type relationDal struct{}

var RelationDal = &relationDal{}

// Relation 数据库用户关系表
type Relation struct {
	gorm.Model
	Id       int64 `gorm:"primaryKey"`
	FollowId int64 // 关注者
	UserId   int64 // 被关注者
}

// CheckRelation 查询用户间关系
func (r *relationDal) CheckRelation(idA int64, idB int64) (bool, error) {
	//查看A是否关注了B
	var db = global.MysqlDB

	var res Relation

	t := db.Where("follow_id = ? AND user_id = ?", idA, idB).First(&res)

	if t.Error != nil {
		return false, t.Error
	} else {
		return true, t.Error
	}
}

// Follow A关注B
func (r *relationDal) Follow(idA int64, idB int64) error {
	//A关注B
	var db = global.MysqlDB

	var res Relation

	res.FollowId = idA
	res.UserId = idB

	relation, _ := r.CheckRelation(idA, idB)

	if relation == true {
		return nil
	} else {
		t := db.Create(&res)
		if t.Error != nil {
			return t.Error
		}
	}
	return nil
}

// Delete A取关B
func (r *relationDal) Delete(idA int64, idB int64) error {
	return global.MysqlDB.Where("follow_id = ? AND user_id = ?", idA, idB).Delete(&Relation{}).Error
}

// GetAllFollow 获取关注用户list
func (r *relationDal) GetAllFollow(id int64) (followIDs []int64, err error) {
	var follows []Relation
	err = global.MysqlDB.Where("follow_id = ?", id).Find(&follows).Error
	if err != nil {
		return
	}
	for _, follow := range follows {
		followIDs = append(followIDs, follow.UserId)
	}
	return
}

// GetAllFollower 获取粉丝用户list
func (r *relationDal) GetAllFollower(id int64) (followerIds []int64, err error) {
	var followers []Relation
	err = global.MysqlDB.Where("user_id = ?", id).Find(&followers).Limit(100).Error
	if err != nil {
		return
	}
	for _, follower := range followers {
		followerIds = append(followerIds, follower.FollowId)
	}
	return
}

// GetAllFriend 获取好友列表
func (r *relationDal) GetAllFriend(Id int64) (friendIds []int64, err error) {
	var followIds []int64
	followIds, err = r.GetAllFollow(Id)
	if err != nil {
		return
	}
	var friends []Relation
	err = global.MysqlDB.Where("user_id IN ? AND follow_id = ?", followIds, Id).Find(&friends).Error
	if err != nil {
		return
	}
	for _, friend := range friends {
		friendIds = append(friendIds, friend.UserId)
	}
	return
}
