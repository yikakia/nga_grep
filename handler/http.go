package handler

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	dots, err := data.GetTimePointsData(start, end, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := []timeseriesResp{}
	for _, v := range dots {
		resp = append(resp, timeseriesResp{
			Timestamp: int64(v.Timestamp * 1000),
			Value:     v.Count,
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
