package handler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/yikakia/nga"
	"github.com/yikakia/nga_grep/internal/observe"
	"github.com/yikakia/nga_grep/model"
	"github.com/yikakia/nga_grep/model/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var syncMeter = sync.OnceValue(func() metric.Meter {
	return otel.Meter("sync")
})

var syncResultCounter = sync.OnceValues(func() (metric.Int64Counter, error) {
	return syncMeter().Int64Counter("sync_result_delta")
})

var syncResultThreadCounter = sync.OnceValues(func() (metric.Int64Counter, error) {
	return syncMeter().Int64Counter("sync_result_delta_thread")
})

func recordSyncResult(delta, deltaThreads int) {
	if counter, err := syncResultCounter(); err == nil {
		counter.Add(context.Background(), int64(delta))
	}
	if counter, err := syncResultThreadCounter(); err == nil {
		counter.Add(context.Background(), int64(deltaThreads))
	}
}

type SyncServerConfig struct {
	Cid string
	Uid string
	Url string
	DB  string
	// 贴子的下限阈值
	ThresholdLow  int
	ThresholdHigh int
	// 下限阈值的倍数
	ThresholdLowFactor float64
	// 上限阈值的倍数
	ThresholdHighFactor float64
	LoopMin             time.Duration
	LoopMax             time.Duration
}

var nextDuration = time.Minute

func updateNextDuration(ctx context.Context, deltaThreads int, cfg SyncServerConfig) {
	tmp := nextDuration
	switch {
	case deltaThreads < cfg.ThresholdLow:
		tmp = time.Duration(float64(nextDuration) * cfg.ThresholdLowFactor)
		tmp -= time.Second
	case deltaThreads > cfg.ThresholdHigh:
		tmp = time.Duration(float64(nextDuration) * cfg.ThresholdHighFactor)
		tmp += time.Second
	}

	nextDuration = durationThreshold(tmp, cfg.LoopMin, cfg.LoopMax)
	slog.InfoContext(ctx, "update next duration", slog.Int("deltaThreads", deltaThreads), slog.Duration("next", nextDuration))
}

// 控制下次调度时间的阈值
// 最小 30s 最多 8min
func durationThreshold(d, loopMin, loopMax time.Duration) time.Duration {
	return max(min(d, loopMax), loopMin)
}

func SyncServer(cfg SyncServerConfig) {
	observe.InitAll()

	c := nga.NewClient(nga.Config{
		BaseUrl:        cfg.Url,
		NgaPassportUid: cfg.Uid,
		NgaPassportCid: cfg.Cid,
	})

	InitDefaultDB(cfg.DB)
	slog.Info("server start success")

	for {
		syncOnce(c, cfg)
		time.Sleep(nextDuration)
	}

}

func syncOnce(c *nga.Client, cfg SyncServerConfig) {
	ctx := context.Background()

	ctx, span := observe.Start(ctx, "sync")
	defer span.End()

	ctx, cspan := observe.Start(ctx, "curl")

	thread, err := c.Thread("706")
	if err != nil {
		slog.ErrorContext(ctx, "query failed.", "err", err.Error())
		cspan.RecordError(err)
		panic(err)
	}
	cspan.End()

	var ts []*model.ThreadLatestData
	for _, t := range thread.Data.T {
		ts = append(ts, &model.ThreadLatestData{
			TID:            t.Tid,
			LastTime:       time.Unix(int64(t.Lastpost), 0),
			LastReplyCount: t.Replies,
		})
	}
	var tids []int
	for _, t := range ts {
		tids = append(tids, t.TID)
	}

	tld := gen.Q.ThreadLatestData

	find, err := tld.WithContext(ctx).Where(tld.TID.In(tids...)).Find()
	if err != nil {
		slog.ErrorContext(ctx, "find from db failed.", "err", err.Error())
		panic(err)
	}

	findMap := lo.SliceToMap(find, func(item *model.ThreadLatestData) (int, *model.ThreadLatestData) {
		return item.TID, item
	})

	var delta, deltaThread int
	for _, t := range ts {
		if v, ok := findMap[t.TID]; ok {
			if v.LastReplyCount != t.LastReplyCount {
				delta += t.LastReplyCount - v.LastReplyCount
				deltaThread++
			}
		} else {
			// 之前没爬，但是这次更新，不计算，仅记录
		}
	}

	ctx, sqlSpan := observe.Start(ctx, "sqlite")
	defer sqlSpan.End()

	err = gen.Q.Transaction(func(tx *gen.Query) error {
		err := tx.WithContext(ctx).ThreadLatestData.Save(ts...)
		if err != nil {
			return fmt.Errorf("insert thread_latest_data failed. err:%w", err)
		}

		err = tx.ThreadCount.WithContext(ctx).Create(&model.ThreadCount{
			DateTime: time.Now().Unix(),
			Count:    delta,
		})
		if err != nil {
			return fmt.Errorf("insert thread_count failed. err:%w", err)
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "update failed.", "err", err.Error())
		panic(err)
	}

	recordSyncResult(delta, deltaThread)

	slog.InfoContext(ctx, "sync success", "delta", delta)
	updateNextDuration(ctx, deltaThread, cfg)

	return
}
