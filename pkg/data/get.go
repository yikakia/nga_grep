package data

import (
	"time"

	"github.com/pkg/errors"
	"github.com/yikakia/nga_grep/model/gen"
)

// 给定一个时间点，以及时间间隔，返回这个时间点所属时间片内的发帖量
// 比如给 2025-04-14 11:33:44 5m, 返回 2025-04-14 11:30:00 到 2025-04-14 11:35:00 内的发帖量 左闭右开
// 最小时间精度是秒
func GetTimePointData(t time.Time, duration time.Duration) (int, error) {
	var posts int

	start := (t.Unix() / int64(duration.Seconds())) * int64(duration.Seconds())
	end := (t.Add(duration).Unix() / int64(duration.Seconds())) * int64(duration.Seconds())

	tc := gen.Q.ThreadCount

	err := tc.Debug().
		Select(tc.Count_.Sum().IfNull(0)).
		Where(tc.DateTime.Gte(start), tc.DateTime.Lt(end)).
		Scan(&posts)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return posts, nil
}

type Dot struct {
	Timestamp int
	Count     int
}

// 给定起始截至时间，以及时间间隔，返回多个时间段的发帖量
func GetTimePointsData(start, end time.Time, duration time.Duration) ([]Dot, error) {
	var dots []Dot
	for cur := start; !cur.After(end); cur = cur.Add(duration) {
		data, err := GetTimePointData(cur, duration)
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
