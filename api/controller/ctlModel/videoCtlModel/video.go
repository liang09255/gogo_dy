package videoCtlModel

import (
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/userCtlModel"
)

// Feed流

type FeedReq struct {
	LatestTime int64 `query:"latest_time"`
}

type FeedResp struct {
	baseCtlModel.BaseResp
	NextTime int64   `json:"next_time"`
	Videos   []Video `json:"video_list"`
}

// 发布视频

type PublishReq struct {
	Title string `form:"title, required" vd:"len($)>0"`
}

type PublishResp struct {
	baseCtlModel.BaseResp
}

// 发布列表

type PublishListReq struct {
	UserID int64 `query:"user_id, required" vd:"$>0"`
}

type PublishListResp struct {
	baseCtlModel.BaseResp
	Videos []Video `json:"video_list"`
}

// Video

type Video struct {
	ID            int64             `json:"id"`
	Author        userCtlModel.User `json:"author"`
	PlayUrl       string            `json:"play_url"`
	CoverUrl      string            `json:"cover_url"`
	FavoriteCount int64             `json:"favorite_count"`
	CommentCount  int64             `json:"comment_count"`
	IsFavorite    bool              `json:"is_favorite"`
	Title         string            `json:"title"`
}
