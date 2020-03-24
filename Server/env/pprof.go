package env

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
)

func MemoryCollector() {

	//初始化一个容器
	memPercent := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memory_percent",
		Help: "memory use percent",
	},
		[]string{"percent"},
	)
	prometheus.MustRegister(memPercent)

	//收集内存使用的百分比
	for {
		//logs.Info("start collect memory used percent!")
		v, err := mem.VirtualMemory()
		if err != nil {
			panic("get memory use percent error: " + err.Error())
		}
		usedPercent := v.UsedPercent
		//logs.Info("get memory use percent:", usedPercent)
		memPercent.WithLabelValues("usedMemory").Set(usedPercent)
		time.Sleep(time.Second * 10)
	}
}
