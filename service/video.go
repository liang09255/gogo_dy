package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	uuid "github.com/satori/go.uuid"
	"main/controller/ctlModel/userCtlModel"
	"main/controller/ctlModel/videoCtlModel"
	"main/dal"
	"main/global"
	"mime/multipart"
	"strconv"
)

type videoService struct{}

var VideoService = &videoService{}

// PublishAction 发布视频服务
func (h *videoService) PublishAction(file *multipart.FileHeader, title string, userID int64) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	uploadFileKey := strconv.FormatInt(userID, 10) + "/" + uuid.NewV4().String() + ".mp4"
	if err := global.AliOSSBucket.PutObject(uploadFileKey, f); err != nil {
		return err
	}
	coverFileKey := uploadFileKey + "?x-oss-process=video/snapshot,t_1000,f_jpg,w_400,h_300,m_fast"

	urlPrefix := "https://" + global.Config.AliOSS.Bucket + ".oss-cn-shenzhen.aliyuncs.com/"
	videoUrl := urlPrefix + uploadFileKey
	coverUrl := urlPrefix + coverFileKey

	err = dal.VideoDal.AddVideo(userID, videoUrl, coverUrl, title)
	if err != nil {
		hlog.Error(err)
	}
	return err
}

// GetPublishList 获取发布的视频列表
func (h *videoService) GetPublishList(userId int64) (ret []videoCtlModel.Video, err error) {
	var videos []dal.Video
	videos, err = dal.VideoDal.GetPublishList(userId)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		var video = videoCtlModel.Video{
			ID:            v.Id,
			Author:        userCtlModel.User{},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
			Title:         v.Title,
		}
		ret = append(ret, video)
	}
	return
}

// Feed 获取视频Feed流
func (h *videoService) Feed(latest int64) (res []videoCtlModel.Video, err error) {
	var videos []dal.Video
	videos, err = dal.VideoDal.GetFeedList(latest)
	if err != nil {
		return nil, err
	}

	for _, v := range videos {
		var video = videoCtlModel.Video{
			ID:            v.Id,
			Author:        userCtlModel.User{},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
			Title:         v.Title,
		}
		res = append(res, video)
	}
	return
}
