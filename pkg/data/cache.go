// 默认使用 gen 包生成的查询器
// 调用前需要设置默认的链接，否则会 panic
package data

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/yikakia/cachalot/core/cache"
	"github.com/yikakia/nga_grep/model/gen"
)

func getTimePointDataWithCache(t time.Time, duration time.Duration) (int, error) {
	get, err := cachalotCache().Get(context.Background(),
		fmt.Sprintf("t:%v_d:%v", t, duration),
		cache.WithOptionCustomField("t", t),
		cache.WithOptionCustomField("duration", duration),
	)
	if err != nil {
		return 0, err
	}

	return get, nil
}

// 给定一个时间点，以及时间间隔，返回这个时间点所属时间片的时间戳
func truncateToDuration(t time.Time, duration time.Duration) int64 {
	return (t.Unix() / int64(duration.Seconds())) * int64(duration.Seconds())
}

// 给定一个时间点，以及时间间隔，返回这个时间点所属时间片内的发帖量
// 比如给 2025-04-14 11:33:44 5m, 返回 2025-04-14 11:30:00 到 2025-04-14 11:35:00 内的发帖量 左闭右开
// 最小时间精度是秒
func getTimePointData(t time.Time, duration time.Duration) (int, error) {
	var posts int

	start := truncateToDuration(t, duration)
	end := truncateToDuration(t.Add(duration), duration)

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
