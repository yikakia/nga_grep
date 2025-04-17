// 默认使用 gen 包生成的查询器
// 调用前需要设置默认的链接，否则会 panic
package data

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yikakia/nga_grep/model/gen"
	"github.com/yikakia/nga_grep/pkg/data/cache"
)

var getDotCache = sync.OnceValue(func() *cache.Cache[int] {
	return cache.NewCache(
		func(args map[string]any) (int, error) {
			t := args["t"].(time.Time)
			duration := args["duration"].(time.Duration)

			data, err := getTimePointData(t, duration)
			if err != nil {
				return 0, err
			}

			return data, nil
		},
		func(args map[string]any) (string, error) {
			t := args["t"].(time.Time)
			duration := args["duration"].(time.Duration)
			k := fmt.Sprintf("t:%v_d:%v", t, duration)
			return k, nil
		})

})

// 给定一个时间点，以及时间间隔，返回这个时间点所属时间片内的发帖量
// 比如给 2025-04-14 11:33:44 5m, 返回 2025-04-14 11:30:00 到 2025-04-14 11:35:00 内的发帖量 左闭右开
// 最小时间精度是秒
func getTimePointData(t time.Time, duration time.Duration) (int, error) {
	var posts int

	start := (t.Unix() / int64(duration.Seconds())) * int64(duration.Seconds())
	end := (t.Add(duration).Unix() / int64(duration.Seconds())) * int64(duration.Seconds())

	tc := gen.Q.ThreadCount

	err := tc.
		Select(tc.Count_.Sum().IfNull(0)).
		Where(tc.DateTime.Gte(start), tc.DateTime.Lt(end)).
		Scan(&posts)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return posts, nil
}

func getTimePointDataWithCache(t time.Time, duration time.Duration) (int, error) {
	v, err := getDotCache().Get(map[string]any{
		"t":        t,
		"duration": duration,
	})
	if err != nil {
		return 0, err
	}

	return v, nil
}
