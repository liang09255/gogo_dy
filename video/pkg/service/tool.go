package service

import (
	"common/ggConv"
	"common/ggIDL/video"
)

func getUserIDFromCommentList(commentList []*video.Comment) []int64 {
	// 获得视频的用户id列表,用一个map存储，重复id可以减少查询
	userIds := make(map[int64]interface{})
	for _, v := range commentList {
		userIds[v.User.Id] = struct{}{}
	}
	return ggConv.Map2Array(userIds)
}

func getUserIDFromVideoList(videoList []*video.Video) []int64 {
	userIds := make(map[int64]interface{})
	for _, v := range videoList {
		userIds[v.Author.Id] = struct{}{}
	}
	return ggConv.Map2Array(userIds)
}
