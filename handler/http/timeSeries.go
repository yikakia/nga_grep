package http

import (
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/conc/pool"
	"github.com/yikakia/nga_grep/pkg/data"
)

func timeSeries(c *gin.Context) {
	var req timeseriesReq

	ctx := c.Request.Context()

	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.WarnContext(ctx, "warn bind query failed.")
		return
	}

	slog.InfoContext(ctx, "bind params succ", req.toAttr())

	start := time.Now().AddDate(0, 0, -1)
	end := time.Now()
	// 默认5分钟
	duration := time.Minute * 5

	if req.StartDate != "" {
		start, err = time.Parse("2006-01-02 15:04", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if req.EndDate != "" {
		end, err = time.Parse("2006-01-02 15:04", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if req.TimeInterval != "" {
		duration, err = time.ParseDuration(req.TimeInterval)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if duration <= 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	if !isAllow(c, start, end, duration) {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	p := pool.New().WithErrors().WithMaxGoroutines(5)
	var dots []timeseriesDots

	p.Go(func() error {
		_dots, err := data.GetTimePointsData(start, end, duration)
		if err != nil {
			return err
		}

		for _, v := range _dots {
			dots = append(dots, timeseriesDots{
				Timestamp: int64(v.Timestamp * 1000),
				Value:     v.Count,
			})
		}

		sort.Slice(dots, func(i, j int) bool {
			return dots[i].Timestamp < dots[j].Timestamp
		})
		return nil
	})

	var apply applyFn
	p.Go(func() error {
		switch req.Indicator {
		case indicatorMA5:
			fn, err := buildMaApplyFn(start, end, duration, 5, func(resps []timeseriesDots, maValues []float64) {
				for i := range resps {
					resps[i].MA5 = maValues[i]
				}
			})
			if err != nil {
				return err
			}
			apply = fn
			return nil
		case indicatorMA10:
			fn, err := buildMaApplyFn(start, end, duration, 10, func(resps []timeseriesDots, maValue []float64) {
				for i := range resps {
					resps[i].MA10 = maValue[i]
				}
			})
			if err != nil {
				return err
			}
			apply = fn
			return nil
		case indicatorBOLL:
			fn, err := buildBollApplyFn(start, end, duration, 20, 2, func(resps []timeseriesDots, bollValues []data.BollingerBand) {
				for i := range resps {
					resps[i].Boll = BollingerBand{
						Middle: bollValues[i].Middle,
						Upper:  bollValues[i].Upper,
						Lower:  bollValues[i].Lower,
					}
				}
			})
			if err != nil {
				return err
			}
			apply = fn
			return nil
		}
		return nil
	})

	err = p.Wait()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if apply != nil {
		apply(dots)

	}

	// 在返回响应前设置缓存头
	c.Header("Cache-Control", "public, max-age=300") // 最多缓存5分钟
	//c.Header("Expires", time.Now().Add(5*time.Minute).UTC().Format(http.TimeFormat))

	c.JSON(http.StatusOK, timeSeriesResp{Data: dots})
}

type applyFn func([]timeseriesDots)

func buildMaApplyFn(start, end time.Time, duration time.Duration, n int, fn func(resps []timeseriesDots, maValues []float64)) (applyFn, error) {
	maN, err := data.GetTimePointsData(start.Add(time.Duration(-n+1)*duration), end, duration)
	if err != nil {
		return nil, err
	}
	maNValues := data.GetMA_N(maN, n)
	return func(resp []timeseriesDots) {
		fn(resp, maNValues)
	}, nil
}

func buildBollApplyFn(start, end time.Time, duration time.Duration, period int, multiplier float64, fn func(resps []timeseriesDots, bollValues []data.BollingerBand)) (applyFn, error) {
	dots, err := data.GetTimePointsData(start.Add(time.Duration(-period+1)*duration), end, duration)
	if err != nil {
		return nil, err
	}
	bolls := data.CalculateBollingerBands(dots, period, multiplier)

	return func(resp []timeseriesDots) {
		fn(resp, bolls)
	}, nil
}
