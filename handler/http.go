package handler

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/model/gen"
	"github.com/yikakia/nga_grep/pkg/data"
)

type RunHttpServerConfig struct {
	Port            string
	CorsAllowOrigin []string
	DB              string
}

func RunHttpServer(cfg RunHttpServerConfig) {
	gen.SetDefault(client.NewDB(cfg.DB))

	r := gin.Default()
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}, // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},      // 允许的头
		ExposeHeaders:    []string{"Content-Length"},                                // 暴露的头
		AllowCredentials: true,                                                      // 允许携带凭证 (例如 cookies)
		AllowOriginFunc: func(origin string) bool { // 允许的源的函数 (更灵活的控制)
			for _, s := range cfg.CorsAllowOrigin {
				if strings.Contains(origin, s) {
					return true
				}
				log.Print(origin)
			}
			log.Printf("hit cors")
			return false
		},
		MaxAge: 12 * time.Hour, // 预检请求的缓存时间
	}

	r.Use(cors.New(config))

	// 监听 /my-path 路径
	r.GET("/api/timeseries", timeSeries)
	err := r.Run(cfg.Port)
	if err != nil {
		panic(err)
	}
}

func timeSeries(c *gin.Context) {
	var req timeseriesReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	}

	p := pool.New().WithErrors().WithMaxGoroutines(5)
	var resp []timeseriesResp

	p.Go(func() error {
		dots, err := data.GetTimePointsData(start, end, duration)
		if err != nil {
			return err
		}

		for _, v := range dots {
			resp = append(resp, timeseriesResp{
				Timestamp: int64(v.Timestamp * 1000),
				Value:     v.Count,
			})
		}

		sort.Slice(resp, func(i, j int) bool {
			return resp[i].Timestamp < resp[j].Timestamp
		})
		return nil
	})

	var apply applyFn
	switch req.Indicator {
	case indicatorMA5:
		p.Go(func() error {
			fn, err := buildMaApplyFn(start, end, duration, 5, func(resps []timeseriesResp, maValues []float64) {
				for i := range resps {
					resps[i].MA5 = maValues[i]
				}
			})
			if err != nil {
				return err
			}
			apply = fn
			return nil
		})
	case indicatorMA10:
		p.Go(func() error {
			fn, err := buildMaApplyFn(start, end, duration, 10, func(resps []timeseriesResp, maValue []float64) {
				for i := range resps {
					resps[i].MA10 = maValue[i]
				}
			})
			if err != nil {
				return err
			}
			apply = fn
			return nil
		})
	}

	err = p.Wait()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if apply != nil {
		apply(resp)

	}

	c.JSON(http.StatusOK, resp)
}

type applyFn func([]timeseriesResp)

func buildMaApplyFn(start, end time.Time, duration time.Duration, n int, fn func(resps []timeseriesResp, maValues []float64)) (applyFn, error) {
	maN, err := data.GetTimePointsData(start.Add(time.Duration(-n+1)*duration), end, duration)
	if err != nil {
		return nil, err
	}
	maNValues := getMA_N(maN, n)
	return func(resp []timeseriesResp) {
		fn(resp, maNValues)
	}, nil
}

func getMA_N(dots []data.Dot, n int) []float64 {
	var maNValues []float64
	values := lo.Map(dots, func(item data.Dot, index int) int {
		return item.Count
	})
	window := lo.Sum(values[:n-1])
	for i, v := range values[n-1:] {
		// i 是第n个
		window += v
		if i > 0 {
			window -= values[i]
		}

		maNValues = append(maNValues, float64(window)/float64(n))
	}
	return maNValues
}

type timeseriesReq struct {
	StartDate    string `form:"startDate"`
	EndDate      string `form:"endDate"`
	TimeInterval string `form:"timeInterval"`
	Indicator    string `form:"indicator"` // 技术指标
}

const (
	indicatorMA5  = "ma5"
	indicatorMA10 = "ma10"
	indicatorBOLL = "boll"
)

type timeseriesResp struct {
	Timestamp int64   `json:"timestamp"` // 毫秒
	Value     int     `json:"value"`
	MA5       float64 `json:"ma5"`
	MA10      float64 `json:"ma10"`
}
