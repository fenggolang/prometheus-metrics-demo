package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var AccessCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_tatol",
	},
	[]string{"method", "path"},
)

var QueueGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "queue_num_total",
	},
	[]string{"name"},
)

var HttpHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_durations_histogram_seconds",
		Help: "http持续时间(单位秒,Histogram--->累积直方图)",
		Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
	},
	[]string{"path"},
)
var HttpDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_durations_seconds",
		Help: "http持续时间(单位秒,Summary--->摘要)",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(AccessCounter,QueueGauge,HttpHistogram,HttpDurations)
}

/**
Prometheus 四大度量指标和应用: https://www.jianshu.com/p/fa5f911003c6

1. Counter (计数器)
2. Gauge (仪表盘)
3. Histogram(累积直方图)
4. Summary（摘要）

腾讯云：使用Prometheus计算百分位数值
https://cloud.tencent.com/developer/news/319419
*/
func main() {

	r := gin.Default()
	// counter 计数器
	r.GET("/counter", func(context *gin.Context) {
		url, _ := url.Parse(context.Request.RequestURI)
		AccessCounter.With(prometheus.Labels{
			"method": context.Request.Method,
			"path":   url.Path,
		}).Add(1)
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// queue_eddycjy 仪表盘
	r.GET("/queue", func(context *gin.Context) {
		numStr := context.Query("num")
		num,_ := strconv.Atoi(numStr)

		QueueGauge.With(prometheus.Labels{
			"name":"queue_eddycjy",
		}).Set(float64(num))
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// http_durations_histogram_seconds 累积直方图
	r.GET("/histogram", func(c *gin.Context) {
		url,_ := url.Parse(c.Request.RequestURI)
		HttpHistogram.With(prometheus.Labels{
			"path":url.Path,
		}).Observe(float64(rand.Intn(30)))
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// summary 摘要
	r.GET("/summary", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDurations.With(prometheus.Labels{
			"path": purl.Path,
		}).Observe(float64(rand.Intn(30)))
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run()
}
