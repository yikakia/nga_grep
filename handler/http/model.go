package http

import (
	"log/slog"
)

type timeseriesReq struct {
	StartDate    string `form:"startDate"`
	EndDate      string `form:"endDate"`
	TimeInterval string `form:"timeInterval"`
	Indicator    string `form:"indicator"` // 技术指标
}

func (t timeseriesReq) toAttr() slog.Attr {
	return slog.Group("timeseriesReq",
		slog.String("startDate", t.StartDate),
		slog.String("endDate", t.EndDate),
		slog.String("timeInterval", t.TimeInterval),
		slog.String("indicator", t.Indicator))
}

const (
	indicatorMA5  = "ma5"
	indicatorMA10 = "ma10"
	indicatorBOLL = "boll"
)

type timeSeriesResp struct {
	Data []timeseriesDots `json:"data"`
}

type timeseriesDots struct {
	Timestamp int64         `json:"timestamp"` // 毫秒
	Value     int           `json:"value"`
	MA5       float64       `json:"ma5"`
	MA10      float64       `json:"ma10"`
	Boll      BollingerBand `json:"boll"`
}
type BollingerBand struct {
	Middle float64 `json:"middle"`
	Upper  float64 `json:"upper"`
	Lower  float64 `json:"lower"`
}
