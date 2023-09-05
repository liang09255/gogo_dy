package repo

import (
	"context"
	"time"
)

// 流程 -> 用户uid对视频vid点赞 (1)-> 缓存存在 -> 修改两步缓存，异步写入relation表(控制点赞的并发)，然后再异步写入video总表里
//	   (2)|										|
//   查无缓存						  修改uid的点赞数,to_uid的获赞数 + 视频vid的点赞数
//      |
//  回数据库查找 -> 这时候数据还是旧的,select for update加锁避免并发修改(否则如果两个数据同时来改),直接修改video总表，以及加上一条点赞记录
//  因为无缓存说明大概率是一些冷门数据，点赞并发量不高，所以可以考虑加锁再改
//  又因为写数据库比写缓存慢，就算多个并发，也能大概率确保，保证同时只有一个进行修改然后返回redis
// 	-> 如果无缓存说明大概率是一些冷门数据，点赞并发量不高 -> 那么就修改数据库

type FavoriteCacheRepo interface {
	GetVideoFavoriteCount(ctx context.Context, vid int64) (int64, bool, error)
	GetUserFavoriteCount(ctx context.Context, uid int64) (int64, bool, error)
	GetUserGetFavoriteCount(ctx context.Context, uid int64) (int64, bool, error)
	SetVideoFavoriteCount(ctx context.Context, vid int64, value int64, expire time.Duration) error
	SetUserFavoriteCount(ctx context.Context, uid int64, value int64, expire time.Duration) error
	SetUseGetFavoriteCount(ctx context.Context, uid int64, value int64, expire time.Duration) error
	IncrUserGetFavoriteCount(ctx context.Context, uid int64) error
	IncrUserFavoriteCount(ctx context.Context, uid int64) error
	IncrVideoFavoriteCount(ctx context.Context, vid int64) error
	DecrUserGetFavoriteCount(ctx context.Context, uid int64) error
	DecrUserFavoriteCount(ctx context.Context, uid int64) error
	DecrVideoFavoriteCount(ctx context.Context, vid int64) error
}
