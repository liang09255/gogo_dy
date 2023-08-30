package repo

import (
	"context"
)

type RelationRepo interface {
	Follow(ctx context.Context, uid int64, targetId int64) error
	Delete(ctx context.Context, uid int64, targetId int64) error
	GetAllFollow(ctx context.Context, id int64) (followIDs []int64, err error)
	GetAllFollower(ctx context.Context, id int64) (followerIds []int64, err error)
	GetAllFriend(ctx context.Context, id int64) (friendIds []int64, err error)
}
