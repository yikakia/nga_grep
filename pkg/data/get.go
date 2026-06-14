package data

import (
	"context"
	"time"

	"github.com/bytedance/gg/gslice"
	"github.com/yikakia/nga_grep/internal/observe"
	"github.com/yikakia/nga_grep/model"
	"github.com/yikakia/nga_grep/model/gen"
)

type Dot struct {
	Timestamp int
	Count     int
}

// 给定起始截至时间，以及时间间隔，返回多个时间段的发帖量
func GetTimePointsData(ctx context.Context, start, end time.Time, duration time.Duration) ([]Dot, error) {
	_, sp := observe.Start(ctx, "GetTimePointsData")
	defer sp.End()

	return getWithSqliteGroupby(ctx, start, end, duration)
}

func getWithCache(start, end time.Time, duration time.Duration) ([]Dot, error) {
	var dots []Dot
	for cur := start; !cur.After(end); cur = cur.Add(duration) {
		data, err := getTimePointDataWithCache(cur, duration)
		if err != nil {
			return nil, err
		}

		dots = append(dots, Dot{
			Timestamp: int((cur.Unix() / int64(duration.Seconds())) * int64(duration.Seconds())),
			Count:     data,
		})
	}
	return dots, nil
}
func getWithSqliteGroupby(ctx context.Context, start, end time.Time, duration time.Duration) ([]Dot, error) {
	startTs := truncateToDuration(start, duration)
	endTs := truncateToDuration(end.Add(duration), duration)
	durationSecond := int64(duration.Seconds())
	tc := gen.Q.ThreadCount

	selectDots, err := tc.WithContext(ctx).SelectDots(startTs, endTs, durationSecond)
	if err != nil {
		return nil, err
	}
	mappedDots := gslice.ToMapValues(selectDots, func(t *model.ThreadCount) int64 {
		return t.DateTime
	})

	var dots []Dot
	for cur := startTs; cur < endTs; cur += durationSecond {
		cnt := 0
		if d, ok := mappedDots[cur]; ok {
			cnt = d.Count
		}
		dots = append(dots, Dot{
			Timestamp: int(cur),
			Count:     cnt,
		})
	}
	return dots, nil
}
