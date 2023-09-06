package domain

import (
	"common/ggConv"
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"context"
	"time"
	"video/internal/dal"
	"video/internal/model"
	"video/internal/mq"
	"video/internal/repo"
)

type FavoriteDomain struct {
	tranRepo        repo.TranRepo
	favoriteRepo    repo.FavoriteRepo
	favoriteCache   repo.FavoriteCacheRepo
	videoRepo       repo.VideoRepo
	videoDetailRepo repo.VideoDetailRepo
}

func NewFavoriteDomain() *FavoriteDomain {
	return &FavoriteDomain{
		tranRepo:        dal.NewTranRepo(),
		favoriteRepo:    dal.NewFavoriteDao(),
		videoRepo:       dal.NewVideoDao(),
		favoriteCache:   dal.NewFavoriteCacheRepo(),
		videoDetailRepo: dal.NewVideoDetailDao(),
	}
}

// FavoriteAction 点赞/取消操作
func (fd *FavoriteDomain) FavoriteAction(ctx context.Context, userid int64, videoid int64, actionType video.ActionType) (err error) {

	// 调用favorite表的获取视频作者id
	toUser, err := fd.videoRepo.GetVideoInfo(ctx, videoid)
	if err != nil {
		ggLog.Error("获得视频作者id错误", err)
		return err
	}
	if actionType == video.ActionType_Add {
		// 如果是点赞操作
		// 1.用户能点赞该视频说明已经把视频和作者相关缓存加载出来了
		// 2. 所以需要判缓存存在的就是点赞者的点赞总数
		// 2-1. 如果缓存存在则缓存增量，异步写入数据库
		// 2-2. 如果缓存不存在，则直接写入点赞记录，异步去写入video的点赞总数
		_, exist, err := fd.favoriteCache.GetUserFavoriteCount(ctx, userid)
		if !exist {
			// 如果用户点赞缓存不存在
			// 直接改数据库,新增记录
			ggLog.Info("缓存不存在")
			// 执行点赞操作
			err = fd.favoriteRepo.PostFavoriteAction(ctx, userid, videoid)
			if err != nil {
				ggLog.Errorf("post favorite action err : %v", err)
				return err
			}
			// 从数据库中获得用户点赞数
			count := fd.favoriteRepo.GetFavoriteCount(ctx, userid)
			// 写入缓存 -- 超时时间需要设置
			err = fd.favoriteCache.SetUserFavoriteCount(ctx, userid, count, 30*time.Second)
			if err != nil && err.Error() != "重复记录" {
				ggLog.Errorf("set user favorite count cache err:%v", err)
				return err
			}
			if err != nil && err.Error() == "重复记录" {
				return nil
			}
			mq.SendVideoFavoriteMsg(mq.FavoriteMessage{Uid: userid, Vid: videoid, Method: int64(actionType)})
		} else {
			// 如果用户缓存存在，则直接更新缓存
			err = fd.favoriteCache.IncrUserFavoriteCount(ctx, userid)
			if err != nil {
				ggLog.Errorf("post favorite action err : %v", err)
				return err
			}
			//  异步写入数据库,避免高并发
			err := mq.AddFavoriteMessage(&mq.FavoriteMessage{
				Vid:    videoid,
				Uid:    userid,
				Method: int64(actionType),
			})
			if err != nil {
				ggLog.Errorf("写入mq错误:%v", err)
				return err
			}
		}

		// 更新作者缓存获赞数
		err = fd.favoriteCache.IncrUserGetFavoriteCount(ctx, toUser.Id)
		if err != nil {
			ggLog.Errorf("incr user favorited count err : %v", err)
			return err
		}
		// 更新视频缓存点赞数
		err = fd.favoriteCache.IncrVideoFavoriteCount(ctx, videoid)
		if err != nil {
			ggLog.Errorf("incr video favorite count err : %v", err)
			return err
		}

	} else {
		// 取消点赞
		_, exist, err := fd.favoriteCache.GetUserFavoriteCount(ctx, userid)
		if !exist {
			// 如果用户点赞缓存不存在
			// 直接改数据库
			err = fd.favoriteRepo.CancelFavoriteAction(ctx, userid, videoid)
			if err != nil {
				ggLog.Errorf("cancel favorite action err : %v", err)
				return err
			}

			// 获得用户点赞数
			count := fd.favoriteRepo.GetFavoriteCount(ctx, userid)
			// 加入缓存 -- 超时时间需要设置
			err = fd.favoriteCache.SetUserFavoriteCount(ctx, userid, count, 10*time.Second)
			if err != nil {
				ggLog.Errorf("set user favorite count cache err:%v", err)
				return err
			}
			mq.SendVideoFavoriteMsg(mq.FavoriteMessage{Uid: userid, Vid: videoid, Method: int64(actionType)})
		} else {
			// 如果用户缓存存在，则直接更新缓存
			err = fd.favoriteCache.DecrUserFavoriteCount(ctx, userid)
			if err != nil {
				ggLog.Errorf("decr favorite action err : %v", err)
				return err
			}
			// 异步写入数据库,避免高并发
			err = mq.AddFavoriteMessage(&mq.FavoriteMessage{
				Uid:    userid,
				Vid:    videoid,
				Method: int64(actionType),
			})
			if err != nil {
				ggLog.Errorf("写入mq错误:%v", err)
				return err
			}
		}
		// 更新作者被赞数
		err = fd.favoriteCache.DecrUserGetFavoriteCount(ctx, toUser.Id)
		if err != nil {
			ggLog.Errorf("decr user favorited count err : %v", err)
			return err
		}
		// 更新视频赞数
		err = fd.favoriteCache.DecrVideoFavoriteCount(ctx, videoid)
		if err != nil {
			ggLog.Errorf("decr video favorite count err : %v", err)
			return err
		}
	}
	return nil
}

// FavoriteList 喜爱列表
func (fd *FavoriteDomain) FavoriteList(ctx context.Context, userid int64) (resp *video.FavoriteListResponse, err error) {
	// 获得视频id列表
	ids, err := fd.favoriteRepo.GetFavoriteList(ctx, userid)
	if err != nil {
		ggLog.Errorf("获得用户:%d 喜爱视频id错误:%v", userid, err)
		return &video.FavoriteListResponse{}, err
	}

	// 根据视频id获得视频列表
	videoList, err := fd.videoRepo.MGetVideoInfo(ctx, ids)
	if err != nil {
		ggLog.Errorf("获得用户:%d 喜爱视频列表错误:%v", userid, err)
		return &video.FavoriteListResponse{}, err
	}

	// 类型转换 - 顺便全部赋值为喜爱视频
	respVideo := fd.videos2Pb(ctx, videoList, true)

	// 初始化，否则会直接panic
	resp = new(video.FavoriteListResponse)
	resp.VideoList = respVideo

	return
}

// 类型转换 - 顺便设置是否为喜爱视频
func (fd *FavoriteDomain) videos2Pb(ctx context.Context, videoList []*model.Video, isFavorite bool) []*video.Video {
	pbs := make([]*video.Video, 0)
	for _, v := range videoList {
		detail := fd.getVideoDetail(ctx, v.Id)
		//获取点赞数和评论数
		p := &video.Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			Title:         v.Title,
			IsFavorite:    isFavorite,
			FavoriteCount: detail.FavoriteCount,
			CommentCount:  detail.CommentCount,
			Author: &user.UserInfoModel{
				Id: v.AuthorId,
			},
		}
		pbs = append(pbs, p)
	}

	return pbs
}

func (fd *FavoriteDomain) getVideoDetail(ctx context.Context, vid int64) *model.VideoDetail {
	detail, err := fd.videoDetailRepo.GetVideoDetail(ctx, vid)
	if err != nil {
		ggLog.Errorf("获得视频:%d 详情错误:%v", vid, err)
	}
	return detail
}

func (fd *FavoriteDomain) GetFavoriteCountByVideoID(ctx context.Context, videoID []int64) map[int64]int64 {
	return fd.favoriteRepo.GetFavoriteCountByVideoID(ctx, videoID)
}

func (fd *FavoriteDomain) CheckFavorite(ctx context.Context, userID int64, videoID []int64) map[int64]bool {
	var res = make(map[int64]bool)
	favoriteIDs, err := fd.favoriteRepo.GetFavoriteList(ctx, userID)
	if err != nil {
		ggLog.Errorf("获得用户:%d 喜爱视频id错误:%v", userID, err)
		return res
	}
	favoriteMap := ggConv.Array2BoolMap(favoriteIDs)
	for _, v := range videoID {
		res[v] = favoriteMap[v]
	}
	return res
}

// 获得点赞数
func (fd *FavoriteDomain) GetFavoriteCount(ctx context.Context, userID int64) int64 {
	// 查询是否有缓存，有缓存直接返回，没缓存就进数据库查询缓存后写回
	count, exist, _ := fd.favoriteCache.GetUserGetFavoriteCount(ctx, userID)
	if !exist {
		vs, err := fd.videoRepo.PublishList(ctx, userID)
		if err != nil {
			ggLog.Errorf("获得用户:%d 喜爱视频id错误:%v", userID, err)
			return 0
		}
		vids := make([]int64, len(vs))
		for _, v := range vs {
			vids = append(vids, v.Id)
		}
		favoriteCountMap := fd.favoriteRepo.GetFavoriteCountByVideoID(ctx, vids)
		var count int64
		for _, v := range favoriteCountMap {
			count += v
		}
		return count
	}
	// 缓存存在则直接返回
	return count
}

// 获得用户的喜欢数
func (fd *FavoriteDomain) GetLikeCount(ctx context.Context, userID int64) int64 {
	count, exist, _ := fd.favoriteCache.GetUserFavoriteCount(ctx, userID)
	if !exist {
		vids, err := fd.favoriteRepo.GetFavoriteList(ctx, userID)
		if err != nil {
			ggLog.Errorf("获得用户:%d 喜爱视频id错误:%v", userID, err)
			return 0
		}
		return int64(len(vids))
	}
	// 缓存存在
	return count
}
