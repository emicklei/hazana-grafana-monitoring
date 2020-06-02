package monitoring

import (
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	monitorInit.Do(initMonitoring)
	before := time.Now().Add(-10 * time.Millisecond)
	for i := 0; i < 102; i++ {
		timerForLabel("test").Update(time.Now().Sub(before))
		time.Sleep(123 * time.Millisecond)
	}
}
