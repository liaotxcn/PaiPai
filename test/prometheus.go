package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// prometheus 测试实例

func main() {
	// 自定义监控指标(项目程序内部)
	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tests",
		Help: "prometheus tests",
	})
	prometheus.MustRegister(temp)

	var i int
	go func() {
		for {
			i++
			if i%2 == 0 {
				temp.Inc()
			}
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe("1234", nil))
}
