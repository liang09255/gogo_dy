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
	relationRepo repo.RelationRepo
	tranRepo     repo.TranRepo
}

func NewRelationDomain() *RelationDomain {
	return &RelationDomain{
		relationRepo: dal.NewRelationDao(),
		tranRepo:     dal.NewTranRepo(),
	}
}

func (ud *RelationDomain) RelationAction(ctx context.Context, userID, targetID int64, actionType int32) error {
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

func (ud *RelationDomain) GetIsFollow(ctx context.Context, myID int64, targetID []int64) map[int64]bool {
	res := make(map[int64]bool)
	defer func() {
		res[myID] = true
	}()
	followIDs, err := ud.relationRepo.GetAllFollow(ctx, myID)
	if err != nil {
		ggLog.Error(err)
		return res
	}
	followIDMap := ggConv.Array2BoolMap(followIDs)
	for _, id := range targetID {
		res[id] = followIDMap[id]
	}
	return res
}
