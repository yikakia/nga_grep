package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/samber/lo"
	"github.com/yikakia/nga"
	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/model"
	"github.com/yikakia/nga_grep/model/gen"
)

type SyncServerConfig struct {
	Cid string
	Uid string
	Url string
	DB  string
}

var nextDuration = time.Second

func updateNextDuration(deltaThreads int) {
	tmp := nextDuration
	switch {
	case deltaThreads < 7:
		tmp = time.Duration(float64(nextDuration) * 1.5)
	case deltaThreads > 14:
		tmp = time.Duration(float64(nextDuration) * 0.7)
	}

	nextDuration = durationThreshold(tmp)
	log.Printf("update next duration. delta:%d next:%v", deltaThreads, nextDuration)
}

// 控制下次调度时间的阈值
// 最小 30s 最多 8min
func durationThreshold(d time.Duration) time.Duration {
	return max(min(d, time.Minute*8), time.Second*30)
}

func SyncServer(cfg SyncServerConfig) {
	c := nga.NewClient(nga.Config{
		BaseUrl:        cfg.Url,
		NgaPassportUid: cfg.Uid,
		NgaPassportCid: cfg.Cid,
	})

	gen.SetDefault(client.NewDB(cfg.DB))

	for {
		sync(c)
		time.Sleep(nextDuration)
	}

}

func sync(c *nga.Client) {
	thread, err := c.Thread("706")
	if err != nil {
		log.Fatal(err)
	}

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

	find, err := tld.Where(tld.TID.In(tids...)).Find()
	if err != nil {
		log.Fatal("sync failed. err:", err)
		return
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

	err = gen.Q.Transaction(func(tx *gen.Query) error {
		err := tx.ThreadLatestData.Save(ts...)
		if err != nil {
			return fmt.Errorf("insert thread_latest_data failed. err:%w", err)
		}

		err = tx.ThreadCount.Create(&model.ThreadCount{
			DateTime: time.Now().Unix(),
			Count:    delta,
		})
		if err != nil {
			return fmt.Errorf("insert thread_count failed. err:%w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal("update failed. err:", err)
	}

	log.Printf("sync success.time:%v delta:%d", time.Now(), delta)
	updateNextDuration(deltaThread)

	return
}
