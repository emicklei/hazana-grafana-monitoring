package monitoring

import (
	"flag"
	"log"
	"net"
	"sync"
	"time"

	graphite "github.com/cyberdelia/go-metrics-graphite"
	"github.com/emicklei/hazana"
	metrics "github.com/rcrowley/go-metrics"
)

var (
	oMonitor       = flag.Bool("m", false, "if true connect to graphite and send metrics")
	oMonitorPrefix = flag.String("p", "hazana", "prefix for metrics")
	timer          metrics.Timer
	gauge          metrics.Gauge
	count          int64
	monitorInit    sync.Once
)

func initMonitoring() {
	log.Println("setup graphite")
	addr, err := net.ResolveTCPAddr("", ":2003")
	if err != nil {
		log.Fatal("ResolveTCPAddr failed ", err)
	}
	go graphite.Graphite(metrics.DefaultRegistry, 1*time.Second, *oMonitorPrefix, addr)
	timer = metrics.NewTimer()
	metrics.Register("call-timer", timer)
	gauge = metrics.NewGauge()
	metrics.Register("goroutines-count", gauge)
}

// Monitored is a Attack decorator that send metrics to graphite
type Monitored struct {
	hazana.Attack
}

// NewMonitor returns a new Monitor decoration on an Attack
func NewMonitor(a hazana.Attack) Monitored {
	return Monitored{a}
}

// Do is part of hazana.Attack
func (m Monitored) Do() hazana.DoResult {
	before := time.Now()
	result := m.Attack.Do()
	if *oMonitor {
		timer.Update(time.Now().Sub(before))
	}
	return result
}

// Setup is part of hazana.Attack
func (m Monitored) Setup(c hazana.Config) error {
	if err := m.Attack.Setup(c); err != nil {
		return err
	}
	if *oMonitor {
		monitorInit.Do(initMonitoring)
	}
	return nil
}

// Clone is part of hazana.Attack
func (m Monitored) Clone() hazana.Attack {
	if *oMonitor {
		count++
		monitorInit.Do(initMonitoring)
		gauge.Update(count)
	}
	return Monitored{m.Attack.Clone()}
}
