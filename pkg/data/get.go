package data

import (
	"time"
)

type Dot struct {
	Timestamp int
	Count     int
}

// 给定起始截至时间，以及时间间隔，返回多个时间段的发帖量
func GetTimePointsData(start, end time.Time, duration time.Duration) ([]Dot, error) {
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
