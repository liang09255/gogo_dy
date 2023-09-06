package dal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/liang09255/lutils/conv"
	"sync"
	"video/internal/database"
	"video/internal/repo"
	"video/pkg/constant"
)

type HotCacheDal struct {
	conn *redis.Client
}

var hotLock sync.RWMutex

var _ repo.HotCacheRepo = (*HotCacheDal)(nil)

func NewHotCacheDal() (hcd *HotCacheDal) {
	return &HotCacheDal{
		conn: database.GetRedis(),
	}
}

func (h *HotCacheDal) AddVidsToHotCache(ctx context.Context, uid int64, vids []int64) error {
	hotLock.RLock()
	defer hotLock.RUnlock()
	// 添加进热门视频列表 set去重
	res := h.conn.SAdd(ctx, constant.HotCacheKey, vids)
	if res.Err() != nil {
		return res.Err()
	}
	// 添加xx访问了xx视频
	for _, vid := range vids {
		cacheKey := getHotUVCacheKey(vid)
		res = h.conn.PFAdd(ctx, cacheKey, uid)
		if res.Err() != nil {
			return res.Err()
		}
	}
	return nil
}

func (h *HotCacheDal) GetAndResetHotCache(ctx context.Context) ([]int64, error) {
	hotLock.Lock()
	defer hotLock.Unlock()
	strRes := h.conn.SMembers(ctx, constant.HotCacheKey)
	if strRes.Err() != nil {
		return nil, strRes.Err()
	}

	vids := make([]int64, 0, len(strRes.Val()))
	for _, vid := range strRes.Val() {
		vids = append(vids, conv.ToInt64(vid))
	}

	res := make([]int64, 0)
	for _, vid := range vids {
		cacheKey := getHotUVCacheKey(vid)
		intRes := h.conn.PFCount(ctx, cacheKey)
		if intRes.Err() != nil {
			return nil, intRes.Err()
		}
		// count大于1000的为热门视频
		if intRes.Val() >= 1000 {
			res = append(res, vid)
		}
		// 重置
		h.conn.Del(ctx, cacheKey)
	}
	h.conn.Del(ctx, constant.HotCacheKey)
	return res, nil
}

func getHotUVCacheKey(vid int64) string {
	return fmt.Sprintf(constant.HotUVCacheKey, vid)
}
