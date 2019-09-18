package monitoring

import (
	"context"
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
	oGraphitePort  = flag.String("g", ":2003", "host:port to connect with Graphite")
	timer          metrics.Timer
	gauge          metrics.Gauge
	count          int64
	monitorInit    sync.Once
)

func initMonitoring() {
	log.Println("[hazana-grafana-monitoring] setup graphite")
	addr, err := net.ResolveTCPAddr("", *oGraphitePort)
	if err != nil {
		log.Fatalf("[hazana-grafana-monitoring] ResolveTCPAddr on [%s] failed error [%v] ", *oGraphitePort, err)
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

// WithMonitor returns a new Monitor decoration on an Attack
func WithMonitor(a hazana.Attack) Monitored {
	return Monitored{a}
}

// Do is part of hazana.Attack
func (m Monitored) Do(ctx context.Context) hazana.DoResult {
	before := time.Now()
	result := m.Attack.Do(ctx)
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
