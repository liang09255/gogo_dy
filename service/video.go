package service

//
//import (
//	"github.com/cloudwego/hertz/pkg/common/hlog"
//	"main/dal"
//	"strconv"
//	"strings"
//)
//
//type DouyinFeedRequest struct {
//	LatestTime *string `json:"latest_time,omitempty"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
//	Token      *string `json:"token,omitempty"`       // 用户登录状态下设置
//}
//
//type DouyinFeedResponse struct {
//	NextTime   *int64          `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
//	StatusCode int64           `json:"status_code"` // 状态码，0-成功，其他值-失败
//	StatusMsg  *string         `json:"status_msg"`  // 返回状态描述
//	VideoList  []dal.VideoInfo `json:"video_list"`  // 视频列表
//}
//
//type DouyinPublishActionResponse struct {
//	StatusCode int64   `json:"status_code"`
//	StatusMsg  *string `json:"status_msg"`
//}
//
//type DouyinPublishActionRequest struct {
//	Token string `json:"token"`
//	Data  []byte `json:"data"`
//	Title string `json:"title"`
//}
//
//type DouyinPublishListRequest struct {
//	Token  string `json:"token"`   // 用户鉴权token
//	UserID string `json:"user_id"` // 用户id
//}
//
//type DouyinPublishListResponse struct {
//	StatusCode int64           `json:"status_code"` // 状态码，0-成功，其他值-失败
//	StatusMsg  *string         `json:"status_msg"`  // 返回状态描述
//	VideoList  []dal.VideoInfo `json:"video_list"`  // 用户发布的视频列表
//}
//
//type videoService struct{}
//
//var VideoService = &videoService{}
//
//// PublishAction 发布视频服务
//func (h *videoService) PublishAction(video dal.Video, token string) error {
//	//根据token获取用户名称
//	s := strings.Split(token, "&")
//	name := s[0]
//
//	var user dal.User
//
//	err := dal.UserDal.GetUserByUserName(name, &user)
//
//	userid := user.ID
//
//	if err != nil {
//		hlog.Error(err)
//		return err
//	}
//
//	video.AuthorId = userid
//	err = dal.VideoDal.AddVideo(video)
//	if err != nil {
//		hlog.Error(err)
//	}
//	return err
//}
//
//// GetPublishList 获取发布的视频列表
//func (h *relationService) GetPublishList(token string, id string, response *[]dal.VideoInfo) interface{} {
//
//	Id, _ := strconv.Atoi(id)
//
//	err := dal.VideoDal.GetPublishList(token, int64(Id), response)
//
//	if err != nil {
//		hlog.Error(err)
//	}
//
//	return err
//}
//
//// Feed 获取视频Feed流
//func (h *videoService) Feed(latest int64, token string, response *[]dal.VideoInfo) error {
//
//	err := dal.VideoDal.GetFeedList(latest, token, response)
//
//	if err != nil {
//		hlog.Error(err)
//	}
//
//	return err
//}
