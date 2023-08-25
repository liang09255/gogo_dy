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
	tranRepo    repo.TranRepo
	commentRepo repo.CommentRepo
	videoRepo   repo.VideoRepo
}

func NewCommentDomain() *CommentDomain {
	return &CommentDomain{
		tranRepo:    dal.NewTranRepo(),
		commentRepo: dal.NewCommentDao(),
		videoRepo:   dal.NewVideoDao(),
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
		commentResp, err := vd.commentRepo.AddComment(ctx, comment)
		if err != nil {
			ggLog.Errorf("用户:%d 发表评论错误:%v", req.UserId, err)
			return nil, err
		}

		// 类型转换
		videoCommentResp = Comment2Pb([]model.Comment{*commentResp})[0]

		// 增加视频表
		err = vd.videoRepo.AddCommentCount(ctx, req.VideoId, 1)
		if err != nil {
			ggLog.Errorf("视频表添加评论数错误，用户:%d 错误:%v", req.UserId, err)
			return nil, err
		}

	} else if req.ActionType == video.ActionType_Cancel {
		comment := &model.Comment{
			ID: req.CommentId,
		}
		err = vd.commentRepo.DeleteComment(ctx, comment)
		if err != nil {
			ggLog.Errorf("用户:%d 删除评论错误:%v", req.UserId, err)
			return nil, err
		}
		// 减少视频表
		err = vd.videoRepo.ReduceCommentCount(ctx, comment.VideoId, 1)
		if err != nil {
			ggLog.Errorf("视频表删除评论数错误，用户:%d 错误:%v", req.UserId, err)
			return nil, err
		}
	}
	conn.Commit()
	return videoCommentResp, nil
}

// CommentList 评论列表
func (vd *CommentDomain) CommentList(ctx context.Context, videoId int64) (commentList []*video.Comment, err error) {
	commentlist, err := vd.commentRepo.GetCommentList(ctx, videoId)
	if err != nil {
		ggLog.Errorf("获得视频:%d 的评论列表错误", videoId)
	}
	// 类型转换
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
