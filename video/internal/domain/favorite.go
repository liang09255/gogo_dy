package domain

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"context"
	"video/internal/dal"
	"video/internal/model"
	"video/internal/repo"
)

type FavoriteDomain struct {
	tranRepo     repo.TranRepo
	favoriteRepo repo.FavoriteRepo
	videoRepo    repo.VideoRepo
}

func NewFavoriteDomain() *FavoriteDomain {
	return &FavoriteDomain{
		tranRepo:     dal.NewTranRepo(),
		favoriteRepo: dal.NewFavoriteDao(),
		videoRepo:    dal.NewVideoDao(),
	}
}

// FavoriteAction 点赞/取消操作
func (fd *FavoriteDomain) FavoriteAction(ctx context.Context, userid int64, videoid int64, actionType video.ActionType) (err error) {
	// 开启事务
	conn := fd.tranRepo.NewTransactionConn()
	conn.Begin()
	defer func() {
		if err != nil {
			conn.Rollback()
		}
	}()
	// 调用favorite表的增加记录
	if actionType == video.ActionType_Add {
		err = fd.favoriteRepo.PostFavoriteAction(ctx, userid, videoid)
		if err != nil {
			return err
		}
		err = fd.videoRepo.AddFavoriteCount(ctx, videoid, 1)
		if err != nil {
			return err
		}
	} else {
		err = fd.favoriteRepo.CancelFavoriteAction(ctx, userid, videoid)
		if err != nil {
			return err
		}

		err = fd.videoRepo.CancelFavoriteCount(ctx, videoid, 1)
		if err != nil {
			return err
		}
	}
	// 走到这一步说明前面全部执行完成,提交事务
	conn.Commit()
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
	respVideo := videos2Pb(videoList, true)

	// 初始化，否则会直接panic
	resp = new(video.FavoriteListResponse)
	//resp.VideoList = make([]*video.Video, 0)

	resp.VideoList = respVideo

	return
}

// 类型转换 - 顺便设置是否为喜爱视频
func videos2Pb(videoList []model.Video, isFavorite bool) []*video.Video {
	pbs := make([]*video.Video, 0)
	for _, v := range videoList {
		p := &video.Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
			Author: &user.UserInfoModel{
				Id: v.AuthorId,
			},
		}
		pbs = append(pbs, p)
	}

	return pbs
}
