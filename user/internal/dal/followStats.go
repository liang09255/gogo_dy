// dal/followStatus.go
package dal

import (
	"common/ggLog"
	"context"
	"github.com/henrylee2cn/goutil/calendar/cron"
	"gorm.io/gorm/clause"

	"time"
	"user/internal/database"
	"user/internal/model"
)

type FollowStatusDal struct {
	conn *database.GormConn
	//reader *kafka.Reader
}

func NewFollowStatusDal() *FollowStatusDal {
	//r := kafka.NewReader("relation-actions") // Initialize kafka reader for topic "relation-actions"
	return &FollowStatusDal{
		conn: database.New(),
		//reader: r,
	}
}

// GetUpdatedRelationsToday 获取当天更新的用户ID
func (fsd *FollowStatusDal) GetUpdatedRelationsToday(ctx context.Context) ([]int64, error) {
	// TODO: 灵活定制更新的时段并解耦
	today := time.Now().Format("2006-01-02")
	var ids []int64
	err := fsd.conn.WithContext(ctx).Model(&model.Relation{}).Where("DATE(updated_at) = ?", today).Pluck("DISTINCT user_id", &ids).Error

	//twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	//currentDate := time.Now().Format("2006-01-02")
	//var ids []int64
	//err := fsd.conn.WithContext(ctx).Model(&model.Relation{}).
	//	Where("DATE(updated_at) BETWEEN ? AND ?", twoDaysAgo, currentDate).
	//	Pluck("DISTINCT user_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetFollowerCountByUserId 获取某个用户的粉丝数量
func (fsd *FollowStatusDal) GetFollowerCountByUserId(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := fsd.conn.WithContext(ctx).Model(&model.Relation{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetFollowCountByUserId 获取某个用户的关注数量
func (fsd *FollowStatusDal) GetFollowCountByUserId(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := fsd.conn.WithContext(ctx).Model(&model.Relation{}).Where("follow_id = ?", userID).Count(&count).Error
	return count, err
}

// SyncFollowStats 同步关注统计
func (fsd *FollowStatusDal) SyncFollowStats(ctx context.Context) error {
	userIds, err := fsd.GetUpdatedRelationsToday(ctx)
	if err != nil {
		return err
	}
	type updateInfo struct {
		UserID         int64
		FollowerCount  int64
		FollowingCount int64
	}

	//var updates []updateInfo // 用来方便调试的，虽然其实没有什么卵用
	//userCache := NewUserCacheRepo()

	for _, userID := range userIds {
		followerCount, err := fsd.GetFollowerCountByUserId(ctx, userID)
		if err != nil {
			// Log and continue with next userID
			ggLog.Errorf("Error fetching follower count for user %d: %v", userID, err)
			continue
		}

		followingCount, err := fsd.GetFollowCountByUserId(ctx, userID)
		if err != nil {
			// Log and continue with next userID
			ggLog.Errorf("Error fetching following count for user %d: %v", userID, err)
			continue
		}

		//// Update Redis cache
		//err = userCache.SetFollowStats(ctx, userID, int(followerCount), int(followingCount), 24*time.Hour)
		//if err != nil {
		//	// handle error
		//	ggLog.Errorf("Error updating Redis cache for user %d: %v", userID, err)
		//	continue
		//}

		// Update the stats table (assuming you have a model for that)
		// Upsert operation
		err = fsd.conn.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"follower_count": followerCount, "following_count": followingCount, "updated_at": time.Now()}),
		}).Create(&model.FollowStats{
			UserID:         userID,
			FollowerCount:  followerCount,
			FollowingCount: followingCount,
			UpdatedAt:      time.Now(),
		}).Error

		if err != nil {
			ggLog.Errorf("Error updating follow stats table for user %d: %v", userID, err)
			continue
		}

		//updates = append(updates, updateInfo{
		//	UserID:         userID,
		//	FollowerCount:  followerCount,
		//	FollowingCount: followingCount,
		//})
	}

	//fmt.Println("我更新了一次")
	//for _, u := range updates {
	//	fmt.Printf("UserID: %d, FollowerCount: %d, FollowingCount: %d\n", u.UserID, u.FollowerCount, u.FollowingCount)
	//}
	return nil
}

func (fsd *FollowStatusDal) StartDailyFollowStatsSync() *cron.Cron {
	c := cron.New()

	// 每天凌晨执行
	c.AddFunc("* * * * *", func() { // 0 0 * * * 是每天凌晨 * * * * *是每分钟更新一次
		// 创建一个新的 context
		ctx := context.Background()
		err := fsd.SyncFollowStats(ctx)
		if err != nil {
			// 这里可以进行错误处理，例如日志记录
			// log.Println("Error syncing follow stats:", err)
			ggLog.Errorf("Error syncing follow stats:", err)
		}
	})

	c.Start()
	ggLog.Info("CRON START")
	return c
}

//func (fsd *FollowStatusDal) ProcessKafkaMessages() {
//	for {
//		message, err := fsd.reader.ReadMessage(context.Background())
//		if err != nil {
//			// Handle error
//			continue
//		}
//		var msgData model.RelationActionMessage
//		err = json.Unmarshal(message.Value, &msgData)
//		if err != nil {
//			// Handle error
//			continue
//		}
//		// Asynchronously update follow stats based on the message content
//		go fsd.UpdateFollowStats(msgData)
//	}
//}
