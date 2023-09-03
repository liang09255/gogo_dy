package domain

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"context"
	"video/internal/dal"
	"video/internal/model"
	"video/internal/repo"
)

type CommentDomain struct {
	tranRepo     repo.TranRepo
	commentRepo  repo.CommentRepo
	videoRepo    repo.VideoRepo
	commentCache repo.CommentCacheRepo
}

func NewCommentDomain() *CommentDomain {
	return &CommentDomain{
		tranRepo:     dal.NewTranRepo(),
		commentRepo:  dal.NewCommentDao(),
		videoRepo:    dal.NewVideoDao(),
		commentCache: dal.NewCommentCacheRepo(),
	}
}

// CommentAction 评论操作
func (vd *CommentDomain) CommentAction(ctx context.Context, req *video.CommentActionRequest) (commentReq *video.Comment, err error) {
	conn := vd.tranRepo.NewTransactionConn()
	conn.Begin()
	defer func() {
		if err != nil {
			conn.Rollback()
		}
	}()

	var videoCommentResp *video.Comment

	// 发表评论
	if req.ActionType == video.ActionType_Add {
		// 构造Comment
		comment := &model.Comment{
			UserId:  req.UserId,
			VideoId: req.VideoId,
			Content: req.CommentText,
		}
		// 缓存视频评论数
		// 查看是否存在，如果存在则直接改，然后插入数据库后，合并计算由定时任务处理，如果不存在则插入数据库

		// TODO 获得评论列表

		commentResp, err := vd.commentRepo.AddComment(ctx, comment)
		if err != nil {
			ggLog.Errorf("用户:%d 发表评论错误:%v", req.UserId, err)
			return nil, err
		}

		// 类型转换
		videoCommentResp = Comment2Pb([]model.Comment{*commentResp})[0]
	} else if req.ActionType == video.ActionType_Cancel {
		comment := &model.Comment{
			ID: req.CommentId,
		}
		err = vd.commentRepo.DeleteComment(ctx, comment)
		if err != nil {
			ggLog.Errorf("用户:%d 删除评论错误:%v", req.UserId, err)
			return nil, err
		}
	}
	conn.Commit()
	return videoCommentResp, nil
}

// CommentList 评论列表
func (vd *CommentDomain) CommentList(ctx context.Context, videoId int64) (commentList []*video.Comment, err error) {
	// 查询缓存中是否存在，不存在了则从数据库中查找
	commentlist, exist, err := vd.commentCache.GetCommentList(ctx, videoId)
	if !exist {
		// 如果不存在，数据库查询
		commentlist, err := vd.commentRepo.GetCommentList(ctx, videoId)
		if err != nil {
			ggLog.Errorf("获得视频:%d 的评论列表错误", videoId)
			return nil, err
		}
		// 类型转换
		commentList = Comment2Pb(commentlist)
		// 写入缓存
		err = vd.commentCache.SetCommentList(ctx, videoId, commentlist)
		if err != nil {
			ggLog.Errorf("写入视频评论缓存错误:%v", err)
		}
		return
	}

	// 如果存在，则直接缓存查询,然后返回
	commentList = Comment2Pb(commentlist)
	return
}

// 类型转换
func Comment2Pb(commentList []model.Comment) []*video.Comment {
	pbs := make([]*video.Comment, 0)
	for _, v := range commentList {
		// 时间格式转换
		// time.Time -> mm-dd
		p := &video.Comment{
			Id: v.ID,
			User: &user.UserInfoModel{
				Id: v.UserId,
			},
			Content:    v.Content,
			CreateDate: v.CreatedAt.Format("01-02"),
		}
		pbs = append(pbs, p)
	}

	return pbs
}

func (vd *CommentDomain) GetCommentCountByVideoID(ctx context.Context, videoIDs []int64) map[int64]int64 {
	return vd.commentRepo.MGetCommentCount(ctx, videoIDs)
}
