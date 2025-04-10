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
	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/model"
	"github.com/yikakia/nga_grep/model/gen"
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

	tc := gen.Q.ThreadCount
	find, err := tc.Where(tc.DateTime.Gte(start.Unix()), tc.DateTime.Lte(end.Unix())).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	divByDuration := lo.GroupBy(find, func(item *model.ThreadCount) int64 {
		return item.DateTime / int64(duration.Seconds()) // 每 duration 为一组
	})

	rMap := lo.MapValues(divByDuration, func(value []*model.ThreadCount, key int64) int {
		return lo.SumBy(value, func(item *model.ThreadCount) int {
			return item.Count
		})
	})
	resp := []timeseriesResp{}

	for cur := start; !cur.After(end); cur = cur.Add(duration) {
		tmpTimestamp := cur.Unix() / int64(duration.Seconds())

		tmpTime := time.Unix(tmpTimestamp*int64(duration.Seconds()), 0)
		retTimeMilli := tmpTime.UnixMilli()
		resp = append(resp, timeseriesResp{
			Timestamp: retTimeMilli,
			Value:     rMap[tmpTimestamp],
		})

	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Timestamp < resp[j].Timestamp
	})

	c.JSON(http.StatusOK, resp)
}

type timeseriesReq struct {
	StartDate    string `form:"startDate"`
	EndDate      string `form:"endDate"`
	TimeInterval string `form:"timeInterval"`
}

type timeseriesResp struct {
	Timestamp int64 `json:"timestamp"` // 毫秒
	Value     int   `json:"value"`
}
