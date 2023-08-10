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
	"time"
)

type videoService struct{}

var VideoService = &videoService{}

// PublishAction 发布视频服务
func (v *videoService) PublishAction(file *multipart.FileHeader, title string, userID int64) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	uploadFileKey := strconv.FormatInt(userID, 10) + "/" + uuid.NewV4().String() + ".mp4"
	if err := global.AliOSSBucket.PutObject(uploadFileKey, f); err != nil {
		return err
	}
	coverFileKey := uploadFileKey + "?x-oss-process=video/snapshot,t_1000,f_jpg,w_720,h_1280,m_fast"

	urlPrefix := "https://" + global.Config.AliOSS.Bucket + ".oss-cn-shenzhen.aliyuncs.com/"
	videoUrl := urlPrefix + uploadFileKey
	coverUrl := urlPrefix + coverFileKey

	err = dal.VideoDal.AddVideo(userID, videoUrl, coverUrl, title)
	if err != nil {
		hlog.Error(err)
	}
	// 添加用户作品数
	err = dal.UserDal.AddWorkCount(userID)
	if err != nil {
		hlog.Error(err)
	}
	return err
}

// GetPublishList 获取发布的视频列表
func (v *videoService) GetPublishList(userId int64) (ret []videoCtlModel.Video, err error) {
	var videos []dal.Video
	videos, err = dal.VideoDal.GetPublishList(userId)
	if err != nil {
		return nil, err
	}

	user, err := dal.UserDal.GetUserInfoById(userId)
	if err != nil {
		return nil, err
	}
	u := userCtlModel.User{
		ID:             user.ID,
		Username:       user.Username,
		FollowCount:    user.FollowerCount,
		FollowerCount:  user.FollowerCount,
		TotalFavorited: user.TotalFavorited,
		WorkCount:      user.WorkCount,
		FavoriteCount:  user.FavoriteCount,
	}

	for _, v := range videos {
		var video = videoCtlModel.Video{
			ID:            v.Id,
			Author:        u,
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
func (v *videoService) Feed(latest, userID int64) (res []videoCtlModel.Video, nextTime time.Time, err error) {
	var videos []dal.Video
	videos, err = dal.VideoDal.GetFeedList(latest)
	if err != nil {
		return
	}

	var authors = make(map[int64]userCtlModel.User)
	var isFavoriteMap = make(map[int64]bool)
	// 登录过的用户
	if userID != 0 {
		authorIds := make([]int64, 0)
		for _, v := range videos {
			authorIds = append(authorIds, v.AuthorId)
		}
		authors, err = UserService.MGetUserInfosMap(authorIds, userID)

		videoIds := make([]int64, 0)
		for _, v := range videos {
			videoIds = append(videoIds, v.Id)
		}
		isFavoriteMap, err = FavoriteService.MGetIsFavorite(videoIds, userID)
	}

	for _, v := range videos {
		var video = videoCtlModel.Video{
			ID:            v.Id,
			Author:        authors[v.AuthorId],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavoriteMap[v.Id],
			Title:         v.Title,
		}
		res = append(res, video)
	}
	if len(videos) != 0 {
		// 最老的视频发布时间为下次获取的开始时间
		nextTime = videos[len(videos)-1].CreatedAt
	} else {
		nextTime = time.Now()
	}
	return
}

func (v *videoService) MGetVideoInfo(videoIds []int64, uid int64) (videos []videoCtlModel.Video, err error) {
	// 获取视频信息
	videoInfos, err := dal.VideoDal.MGetVideoInfo(videoIds)
	if err != nil {
		return nil, err
	}

	// 获取作者信息
	authorIds := make([]int64, 0)
	for _, v := range videoInfos {
		authorIds = append(authorIds, v.AuthorId)
	}
	authors, err := UserService.MGetUserInfo(authorIds, uid)
	if err != nil {
		hlog.Error(err)
	}
	authorMap := make(map[int64]userCtlModel.User)
	for _, v := range authors {
		authorMap[v.ID] = v
	}

	// 获取用户是否收藏
	isFavoriteMap, err := FavoriteService.MGetIsFavorite(videoIds, uid)
	if err != nil {
		hlog.Error(err)
	}

	for _, v := range videoInfos {
		var video = videoCtlModel.Video{
			ID:            v.Id,
			Author:        authorMap[v.AuthorId],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavoriteMap[v.Id],
			Title:         v.Title,
		}
		videos = append(videos, video)
	}
	return
}

func (v *videoService) AddFavoriteCount(videoID int64) error {
	return dal.VideoDal.AddFavoriteCount(videoID)
}

func (v *videoService) AddCommentCount(videoID int64) error {
	return dal.VideoDal.AddCommentCount(videoID)
}

func (v *videoService) ReduceFavoriteCount(videoID int64) error {
	return dal.VideoDal.ReduceFavoriteCount(videoID)
}

func (v *videoService) ReduceCommentCount(videoID int64) error {
	return dal.VideoDal.ReduceCommentCount(videoID)
}
