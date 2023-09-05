package domain

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"context"
	"time"
	"video/internal/dal"
	"video/internal/model"
	"video/internal/repo"
)

type VideoDomain struct {
	videoRepo repo.VideoRepo
	tranRepo  repo.TranRepo
}

func NewVideoDomain() *VideoDomain {
	return &VideoDomain{
		videoRepo: dal.NewVideoDao(),
		tranRepo:  dal.NewTranRepo(),
	}
}

func (vd *VideoDomain) GetVideo(ctx context.Context, videoid int64) (*model.Video, error) {
	return vd.videoRepo.GetVideoInfo(ctx, videoid)
}

func (vd *VideoDomain) NewVideo(ctx context.Context, userID int64, videoURL, coverURL, title string) error {
	return vd.videoRepo.NewVideo(ctx, userID, videoURL, coverURL, title)
}

func (vd *VideoDomain) PublishList(ctx context.Context, userID int64) ([]*video.Video, error) {
	vs, err := vd.videoRepo.PublishList(ctx, userID)
	return videoModel2pbVideo(vs), err
}

func (vd *VideoDomain) Feed(ctx context.Context, latestTime int64) ([]*video.Video, int64, error) {
	vs, err := vd.videoRepo.Feed(ctx, latestTime)
	if err != nil {
		return nil, 0, err
	}
	var nextTime int64
	for _, v := range vs {
		t := v.CreatedAt.UnixMicro()
		if t > nextTime {
			nextTime = t
		}
	}
	videos := videoModel2pbVideo(vs)
	if len(videos) != 0 {
		// 最老的视频发布时间为下次获取的开始时间
		nextTime = vs[len(vs)-1].CreatedAt.UnixMicro()
	} else {
		nextTime = time.Now().UnixMicro()
	}
	return videos, nextTime, err
}

func videoModel2pbVideo(vs []*model.Video) []*video.Video {
	if vs == nil || len(vs) == 0 {
		return nil
	}
	var videos []*video.Video
	for _, v := range vs {
		pbVideo := &video.Video{
			Id:       v.Id,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			//FavoriteCount: v.FavoriteCount,
			//CommentCount:  v.CommentCount,
			Title:  v.Title,
			Author: &user.UserInfoModel{Id: v.AuthorId},
		}
		videos = append(videos, pbVideo)
	}
	return videos
}

func (vd *VideoDomain) GetTotalVideoCount(ctx context.Context, userID int64) int64 {
	vids, err := vd.videoRepo.PublishList(ctx, userID)
	if err != nil {
		ggLog.Error("获取视频数量失败:", err)
		return 0
	}
	return int64(len(vids))
}
