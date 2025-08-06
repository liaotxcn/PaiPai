package main

import (
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/prometheus"
	"time"
)

func main() {
	pCfg := prometheus.Config{
		Host: "0.0.0.0",
		Port: 8080,
		Path: "/metrics",
	}
	prometheus.StartAgent(pCfg)

	gaugeVec := metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: "zeromicro",
		Subsystem: "tests",
		Name:      "go_zero_test",
		Help:      "go_zero prometheus metric",
		Labels:    []string{"path"},
	})

	var i int
	for {
		i++
		if i%2 == 0 {
			gaugeVec.Inc("/user")
		}
		time.Sleep(time.Second)
	}
}
