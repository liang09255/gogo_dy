package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlModel/videoCtlModel"
	"main/dal"
	"main/utils/conv"
)

type favoriteService struct{}

var FavoriteService = &favoriteService{}

func (h *favoriteService) GetFavoriteList(ctx context.Context, userId int64) (videos []videoCtlModel.Video, err error) {
	var videoIds []int64
	videoIds, err = dal.FavoriteDal.GetFavoriteList(ctx, userId)
	if err != nil {
		return nil, err
	}
	videoInfos, err := VideoService.MGetVideoInfo(videoIds, userId)
	for _, video := range videoInfos {
		videos = append(videos, video)
	}
	return
}

func (h *favoriteService) PostFavoriteAction(ctx context.Context, userId int64, videoId int64, actionType int32) (string, error) {
	var err error
	msg := "没有指定的类型"
	if actionType == 1 {
		err = dal.FavoriteDal.PostFavoriteAction(ctx, userId, videoId)
		if err := VideoService.AddFavoriteCount(videoId); err != nil {
			hlog.Error(err)
		}
		msg = "点赞成功"
	} else if actionType == 2 {
		err = dal.FavoriteDal.CancelFavoriteAction(ctx, userId, videoId)
		if err := VideoService.ReduceFavoriteCount(videoId); err != nil {
			hlog.Error(err)
		}
		msg = "取消点赞成功"
	}
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (h *favoriteService) MGetIsFavorite(videoIds []int64, uid int64) (ret map[int64]bool, err error) {
	favoriteIds, err := dal.FavoriteDal.GetFavoriteList(context.Background(), uid)
	if err != nil {
		return nil, err
	}
	ret = make(map[int64]bool)
	videoIdMap := conv.Array2Map(videoIds)
	for _, favoriteId := range favoriteIds {
		if _, ok := videoIdMap[favoriteId]; ok {
			ret[favoriteId] = true
		} else {
			ret[favoriteId] = false
		}
	}
	return
}
