package dal

import (
	"main/global"
)

// Relation 数据库用户关系表
type Relation struct {
	Id           int64 `gorm:"primaryKey"`
	User         User  `gorm:"foreignKey:FollowId"` // 关注人
	FollowId     int64
	FollowedUser User `gorm:"foreignKey:UserId"` // 被关注人
	UserId       int64
	CreateTime   int64 `gorm:"autoCreateTime"`
}
type relationDal struct{}

var RelationDal = &relationDal{}

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

	err := r.AddFollowCount(idA)

	if err != nil {
		return err
	}
	err = r.AddFollowerCount(idB)

	return err
}

// Delete A取关B
func (r *relationDal) Delete(idA int64, idB int64) error {
	//A取消关注B
	var db = global.MysqlDB

	var res Relation

	t := db.Where("follow_id = ? AND user_id = ?", idA, idB).Delete(&res)

	if t.Error != nil {
		return t.Error
	}

	err := r.SubFollowCount(idA)

	if err != nil {
		return err
	}
	err = r.SubFollowerCount(idB)

	return err
}

// AddFollowCount 增加用户关注数
func (r *relationDal) AddFollowCount(id int64) error {
	var db = global.MysqlDB
	var user User
	t := db.Find(&user, id)

	if t.Error != nil {
		return t.Error
	}

	user.FollowCount += 1

	t = db.Save(&user)

	return t.Error
}

// AddFollowerCount 增加用户粉丝数
func (r *relationDal) AddFollowerCount(id int64) error {
	var db = global.MysqlDB
	var user User
	t := db.Find(&user, id)

	if t.Error != nil {
		return t.Error
	}

	user.FollowerCount += 1

	t = db.Save(&user)

	return t.Error
}

// SubFollowCount 减少用户关注数
func (r *relationDal) SubFollowCount(id int64) error {
	var db = global.MysqlDB
	var user User
	t := db.Find(&user, id)

	if t.Error != nil {
		return t.Error
	}

	user.FollowCount -= 1

	t = db.Save(&user)

	return t.Error
}

// SubFollowerCount 增加用户关注数
func (r *relationDal) SubFollowerCount(id int64) error {
	var db = global.MysqlDB
	var user User
	t := db.Find(&user, id)

	if t.Error != nil {
		return t.Error
	}

	user.FollowerCount -= 1

	t = db.Save(&user)

	return t.Error
}

// GetAllFollow 获取关注用户list
func (r *relationDal) GetAllFollow(token string, id int64, resp *[]UserInfoResponse) error {
	var db = global.MysqlDB
	res := make([]Relation, 1000)

	db.Where("follow_id = ?", id).Find(&res)
	var tmp UserInfoResponse
	for k := range res {
		id := res[k].UserId

		err := UserDal.GetUserInfoById(id, &tmp)
		if err != nil {
			return err
		}

		//todo 用户表缺少字段
		//tmp.IsFollow = true
		*resp = append(*resp, tmp)
	}
	return nil
}

// GetAllFollower 获取粉丝用户list
func (r *relationDal) GetAllFollower(token string, id int64, resp *[]UserInfoResponse) error {
	var db = global.MysqlDB
	res := make([]Relation, 1000)

	db.Where("user_id = ?", id).Find(&res)
	var tmp UserInfoResponse
	for k := range res {
		id := res[k].FollowId

		err := UserDal.GetUserInfoById(id, &tmp)
		if err != nil {
			return err
		}

		*resp = append(*resp, tmp)
	}
	return nil
}

// GetAllFriend 获取好友列表
func (r *relationDal) GetAllFriend(token string, Id int64, resp *[]UserInfoResponse) error {
	var db = global.MysqlDB
	res := make([]Relation, 1000)

	db.Where("user_id = ?", Id).Find(&res)
	var tmp UserInfoResponse
	var temp Relation
	for k := range res {
		id := res[k].FollowId
		//粉丝id
		err := UserDal.GetUserInfoById(id, &tmp)
		if err != nil {
			return err
		}
		//判断是否互相关注
		t := db.Where("user_id = ? AND follow_id = ?", id, Id).First(&temp)
		if t.Error == nil {
			*resp = append(*resp, tmp)
			//二者互相关注，为朋友关系
		}
	}
	return nil
}
