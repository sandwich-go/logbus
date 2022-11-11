package monitor

import (
	"math"
	"testing"
	"time"
)

func TestGetProcessStat(t *testing.T) {
	ps := GetProcessStat(time.Second)
	t.Log(ps.CPUUsage(), ps.ResidentMemory())
	go task(1)
	time.Sleep(time.Second)
	t.Log(ps.CPUUsage(), ps.ResidentMemory())
	time.Sleep(time.Second * 2)
	t.Log(ps.CPUUsage(), ps.ResidentMemory())
	time.Sleep(time.Second * 3)
	t.Log(ps.CPUUsage(), ps.ResidentMemory())
}

func task(id int) {
	var j float64 = 0.0
	var step float64 = 0.1
	for j = 0.0; j < 8*2*math.Pi; j += step {
		compute(1000.0, math.Sin(j)/2.0+0.5, id)
	}
}

func compute(t, percent float64, id int) {
	// t 总时间，转换为纳秒
	var r int64 = 1000 * 1000
	totalNanoTime := t * (float64)(r)               // 纳秒
	runtime := totalNanoTime * percent              // 纳秒
	sleeptime := totalNanoTime - runtime            // 纳秒
	starttime := time.Now().UnixNano()              // 当前的纳秒数
	d := time.Duration(sleeptime) * time.Nanosecond // 休眠时间
	for float64(time.Now().UnixNano())-float64(starttime) < runtime {
		// 此处易出错：只能用UnixNano而不能使用Now().Unix()
		// 因为Unix()的单位是秒，而整个运行周期
	}
	time.Sleep(d)
}
