package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/dal"
	"main/model"
	"strings"
)

type videoService struct{}

var VideoService = &videoService{}

// PublishAction 发布视频服务
func (h *videoService) PublishAction(video model.Video, token string) error {
	//根据token获取用户名称
	s := strings.Split(token, "&")
	name := s[0]

	userid, err := dal.UserDal.GetUserIDbyname(name)

	if err != nil {
		hlog.Error(err)
		return err
	}

	video.AuthorId = userid
	err = dal.VideoDal.AddVideo(video)
	if err != nil {
		hlog.Error(err)
	}
	return err
}

// GetPublishList 获取发布的视频列表
func (h *relationService) GetPublishList(token string, id string, response *[]model.VideoInfo) interface{} {
	Id := GetInt64Bystring(id)

	err := dal.VideoDal.GetPublishList(token, Id, response)

	if err != nil {
		hlog.Error(err)
	}

	return err
}

// Feed 获取视频Feed流
func (h *videoService) Feed(latest int64, token string, response *[]model.VideoInfo) error {

	err := dal.VideoDal.GetFeedList(latest, token, response)

	if err != nil {
		hlog.Error(err)
	}

	return err
}
