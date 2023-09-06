package domain

import (
	"common/ggLog"
	"context"
	"time"
	"video/internal/dal"
	"video/internal/repo"
)

type HotDomain struct {
	ticker   *time.Ticker
	hotCache repo.HotCacheRepo
}

func NewHotDomain() (hd *HotDomain) {
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		hd.Start()
	}()
	return &HotDomain{
		ticker:   ticker,
		hotCache: dal.NewHotCacheDal(),
	}
}

func (hd *HotDomain) HotStatistics(uid int64, vids []int64) {
	err := hd.hotCache.AddVidsToHotCache(context.Background(), uid, vids)
	if err != nil {
		ggLog.Error("HotStatistics", "AddVidsToHotCache", err)
	}
}

func (hd *HotDomain) Start() {
	go func() {
		for {
			select {
			case <-hd.ticker.C:
				hd.sendHotStatisticsMsg()
			}
		}
	}()
}

func (hd *HotDomain) sendHotStatisticsMsg() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	vids, err := hd.hotCache.GetAndResetHotCache(ctx)
	if err != nil {
		ggLog.Error("sendHotStatisticsMsg", "GetAndResetHotCache", err)
	}
	if len(vids) == 0 {
		return
	}
	//TODO 后续完善
}
