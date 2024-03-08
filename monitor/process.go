package monitor

import (
	"math"
	"os"
	"sync"
	"time"

	"github.com/prometheus/procfs"
)

var instance = newStat()
var minInterval = time.Second * 15

type ProcessStata interface {
	CPUUsage() float64
	ResidentMemory() int
}

// GetProcessStat 获取进程状态查看器
// interval: 控制获取状态的频率，默认30s，最小15s，
func GetProcessStat(interval time.Duration) ProcessStata {
	if minInterval <= interval && instance.statInterval > interval {
		instance.statInterval = interval
	}
	return instance
}

func newStat() *processStata {
	ps := &processStata{
		statInterval: minInterval * 2,
	}
	pid := os.Getpid()
	var err error
	ps.proc, err = procfs.NewProc(pid)
	if err != nil {
		ps.canStat = false
	} else {
		ps.canStat = true
	}
	ps.stat()
	return ps
}

type processStata struct {
	sync.Mutex
	proc               procfs.Proc
	lastStat, prevStat procfs.ProcStat
	lastTime, prevTime time.Time
	statInterval       time.Duration
	canStat            bool
}

func (p *processStata) CPUUsage() float64 {
	if !p.canStat {
		return 0
	}
	p.stat()

	if p.lastTime.Unix() == p.prevTime.Unix() || p.prevTime.IsZero() {
		return 0
	}
	return math.Trunc(((p.lastStat.CPUTime()-p.prevStat.CPUTime())/float64(p.lastTime.Unix()-p.prevTime.Unix()))*10000) / 10000
}

func (p *processStata) ResidentMemory() int {
	if !p.canStat {
		return 0
	}
	p.stat()

	return p.lastStat.ResidentMemory()
}

func (p *processStata) stat() {
	now := time.Now()
	needSync := now.After(p.lastTime.Add(p.statInterval)) || p.lastTime.IsZero() || p.prevTime.IsZero()
	if !needSync {
		return
	}
	p.Lock()
	defer p.Unlock()
	if !p.lastTime.IsZero() && !p.prevTime.IsZero() && now.Before(p.lastTime.Add(p.statInterval)) {
		return
	}
	p.prevStat = p.lastStat
	p.prevTime = p.lastTime
	var err error
	p.lastStat, err = p.proc.Stat()
	if err != nil {
		return
	}
	p.lastTime = now
}
