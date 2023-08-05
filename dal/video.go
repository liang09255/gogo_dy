package dal

import (
	"main/global"
	"main/model"
	"strconv"
)

type videoDal struct{}

var VideoDal = &videoDal{}

// AddVideo 上传视频
func (h *videoDal) AddVideo(video model.Video) error {
	var db = global.MysqlDB
	//创建视频记录
	t := db.Create(&video)

	return t.Error
}

// GetPublishList 获取某个用户发表过的视频列表
func (h *videoDal) GetPublishList(token string, id int64, resp *[]model.VideoInfo) error {

	var db = global.MysqlDB
	res := make([]model.Video, 1000)

	t := db.Where("author_id = ?", id).Find(&res)

	if t.Error != nil {
		return t.Error
	}

	var userinfo model.Userinfo

	err := UserDal.GetUserInfo(token, strconv.FormatInt(id, 10), &userinfo)

	//判断该视频用户是否是自己关注过的
	fg, err := RelationDal.CheckRelation(id, userinfo.ID)

	if fg == true {
		userinfo.IsFollow = true
	} else {
		userinfo.IsFollow = false
	}

	if err != nil {
		return err
	}
	var tmp model.VideoInfo
	for k := range res {
		tmp.Author = userinfo
		tmp.ID = res[k].Id
		tmp.CommentCount = res[k].CommentCount
		tmp.FavoriteCount = res[k].FavoriteCount
		tmp.Title = res[k].Title
		tmp.PlayURL = res[k].PlayUrl
		tmp.CoverURL = res[k].CoverUrl

		//是否是点赞过的作品
		tmp.IsFavorite = true
		*resp = append(*resp, tmp)
	}

	return nil
}

// GetFeedList 获取视频feed列表
func (h *videoDal) GetFeedList(latest int64, token string, response *[]model.VideoInfo) error {

	var db = global.MysqlDB

	res := make([]model.Video, 1000)

	t := db.Order("create_time DESC").Find(&res).Limit(30)

	if t.Error != nil {
		return t.Error
	}
	var userinfo model.Userinfo
	var tmp model.VideoInfo
	for k := range res {
		//获取视频作者用户信息
		err := UserDal.GetUserInfo(token, strconv.FormatInt(res[k].AuthorId, 10), &userinfo)
		if err != nil {
			return err
		}
		tmp.Author = userinfo
		tmp.ID = res[k].Id
		tmp.CommentCount = res[k].CommentCount
		tmp.FavoriteCount = res[k].FavoriteCount
		tmp.Title = res[k].Title
		//用户播放视频的url
		tmp.PlayURL = "http://192.168.2.59:8080/douyin/find?play_url=" + res[k].PlayUrl
		tmp.CoverURL = "http://192.168.2.59:8080/douyin/find?play_url=" + res[k].CoverUrl

		//判断用户是否点赞了这个视频 todo
		tmp.IsFavorite = false

		*response = append(*response, tmp)
	}

	return nil
}
