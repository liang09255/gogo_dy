package dal

import (
	"main/global"
)

// Video 视频
type Video struct {
	Id            int64  `gorm:"primaryKey"`
	AuthorId      int64  `gorm:"not null"`
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64
	CommentCount  int64
	Title         string `gorm:"not null"`
	CreateTime    int64  `gorm:"autoCreateTime"`
}

// VideoInfo 视频信息
type VideoInfo struct {
	Author        UserInfoResponse `json:"author"`         // 视频作者信息
	CommentCount  int64            `json:"comment_count"`  // 视频的评论总数
	CoverURL      string           `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64            `json:"favorite_count"` // 视频的点赞总数
	ID            int64            `json:"id"`             // 视频唯一标识
	IsFavorite    bool             `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string           `json:"play_url"`       // 视频播放地址
	Title         string           `json:"title"`          // 视频标题
}

type videoDal struct{}

var VideoDal = &videoDal{}

// AddVideo 上传视频
func (h *videoDal) AddVideo(video Video) error {
	var db = global.MysqlDB
	//创建视频记录
	t := db.Create(&video)

	return t.Error
}

// GetPublishList 获取某个用户发表过的视频列表
func (h *videoDal) GetPublishList(token string, id int64, resp *[]VideoInfo) error {

	var db = global.MysqlDB
	res := make([]Video, 1000)

	t := db.Where("author_id = ?", id).Find(&res)

	if t.Error != nil {
		return t.Error
	}

	var user UserInfoResponse

	err := UserDal.GetUserInfoById(id, &user)
	//err := UserDal.GetUserInfo(token, strconv.FormatInt(id, 10), &userinfo)
	//todo 用户信息字段缺少

	//判断该视频用户是否是自己关注过的
	//fg, err := RelationDal.CheckRelation(id, userinfo.ID)

	//if fg == true {
	//	userinfo.IsFollow = true
	//} else {
	//	userinfo.IsFollow = false
	//}

	if err != nil {
		return err
	}
	var tmp VideoInfo
	for k := range res {
		tmp.Author = user
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
func (h *videoDal) GetFeedList(latest int64, token string, response *[]VideoInfo) error {

	var db = global.MysqlDB

	res := make([]Video, 1000)

	t := db.Order("create_time DESC").Find(&res).Limit(30)

	if t.Error != nil {
		return t.Error
	}
	var userinfo UserInfoResponse
	var tmp VideoInfo
	for k := range res {
		//获取视频作者用户信息
		err := UserDal.GetUserInfoById(res[k].AuthorId, &userinfo)
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
