package model

import (
	"gorm.io/gen"
)

type Query interface {
	// SELECT
	//   sum(count) as count,
	//   (date_time / @durationSecond) * @durationSecond as date_time
	//  FROM @@table
	//  WHERE date_time >= @start and date_time < @end
	//  GROUP BY
	//   date_time / @durationSecond
	SelectDots(start, end, durationSecond int64) ([]*gen.T, error)
}
