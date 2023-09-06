package domain

import (
	"common/ggConv"
	"common/ggLog"
	"context"
	"fmt"
	"user/internal/dal"
	"user/internal/repo"
)

const (
	FollowAction   = 1
	UnFollowAction = 2
)

type RelationDomain struct {
	relationRepo  repo.RelationRepo
	relationCache repo.RelationCacheRepo
	tranRepo      repo.TranRepo
}

func NewRelationDomain() *RelationDomain {
	return &RelationDomain{
		relationRepo:  dal.NewRelationDao(),
		relationCache: dal.NewRelationCacheDal(),
		tranRepo:      dal.NewTranRepo(),
	}
}

func (ud *RelationDomain) RelationAction(ctx context.Context, userID, targetID int64, actionType int32) error {
	defer func() {
		ud.relationCache.Delete(ctx, userID)
	}()
	if actionType == FollowAction {
		return ud.relationRepo.Follow(ctx, userID, targetID)
	} else if actionType == UnFollowAction {
		return ud.relationRepo.Delete(ctx, userID, targetID)
	}
	return fmt.Errorf("invalid action type, action_type: %d", actionType)
}

func (ud *RelationDomain) GetFollowList(ctx context.Context, userid int64) (users []int64, err error) {
	return ud.relationRepo.GetAllFollow(ctx, userid)
}

func (ud *RelationDomain) GetFollowerList(ctx context.Context, userid int64) (users []int64, err error) {
	return ud.relationRepo.GetAllFollower(ctx, userid)
}

func (ud *RelationDomain) GetFriendList(ctx context.Context, userid int64) (users []int64, err error) {
	return ud.relationRepo.GetAllFriend(ctx, userid)
}

func (ud *RelationDomain) GetIsFollow(ctx context.Context, myID int64, targetID []int64) (ret map[int64]bool) {
	ret = make(map[int64]bool)
	defer func() {
		ret[myID] = true
	}()

	// 查询是否存在缓存
	if ud.relationCache.KeyExist(ctx, myID) {
		// 查询是否关注
		res, err := ud.relationCache.MEXISTS(ctx, myID, targetID)
		if err != nil {
			ggLog.Error(err)
		}

		for k, id := range targetID {
			// 绝大部分请求都会为false
			if res[k] == false {
				ret[id] = res[k]
				continue
			}
			// 为true需要再去数据库确认一遍
			isFollow, err := ud.relationRepo.IsFollow(ctx, myID, id)
			if err != nil {
				ggLog.Error(err)
				ret[id] = isFollow
				continue
			}
			ret[id] = isFollow
			// 如果数据库数据和缓存不一致，删除缓存待后续重建
			if isFollow != res[k] {
				ud.relationCache.Delete(ctx, myID)
			}
		}
	}
	// 不存在缓存，查询数据库
	followIDs, err := ud.relationRepo.GetAllFollow(ctx, myID)
	if err != nil {
		ggLog.Error(err)
		return ret
	}
	followIDMap := ggConv.Array2BoolMap(followIDs)
	for _, id := range targetID {
		ret[id] = followIDMap[id]
	}
	// 写缓存
	err = ud.relationCache.MADD(ctx, myID, targetID)
	if err != nil {
		ggLog.Error(err)
	}
	return ret
}
