// Copyright (c) 2018 SpiralScout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package grpc

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	rrpc "github.com/spiral/php-grpc"
	"github.com/spiral/roadrunner"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/service/metrics"
	"github.com/spiral/roadrunner/util"
	"google.golang.org/grpc/status"
)

func init() {
	cobra.OnInitialize(func() {
		svc, _ := rr.Container.Get(metrics.ID)
		mtr, ok := svc.(*metrics.Service)
		if !ok || !mtr.Enabled() {
			return
		}

		ht, _ := rr.Container.Get(rrpc.ID)
		if gp, ok := ht.(*rrpc.Service); ok {
			collector := newCollector()

			// register metrics
			mtr.MustRegister(collector.callCounter)
			mtr.MustRegister(collector.callDuration)
			mtr.MustRegister(collector.workersMemory)

			// collect events
			gp.AddListener(collector.listener)

			// update memory usage every 10 seconds
			go collector.collectMemory(gp, time.Second*10)
		}
	})
}

// listener provide debug callback for system events. With colors!
type metricCollector struct {
	callCounter   *prometheus.CounterVec
	callDuration  *prometheus.HistogramVec
	workersMemory prometheus.Gauge
}

func newCollector() *metricCollector {
	return &metricCollector{
		callCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rr_grpc_call_total",
				Help: "Total number of handled grpc requests after server restart.",
			},
			[]string{"code"},
		),
		callDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "rr_grpc_call_duration_seconds",
				Help: "GRPC call duration.",
			},
			[]string{"code"},
		),
		workersMemory: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "rr_grpc_workers_memory_bytes",
				Help: "Memory usage by GRPC workers.",
			},
		),
	}
}

// listener listens to http events and generates nice looking output.
func (c *metricCollector) listener(event int, ctx interface{}) {
	if event == rrpc.EventUnaryCall {
		uc := ctx.(*rrpc.UnaryCallEvent)
		code := "Unknown"

		if uc.Error == nil {
			code = "Ok"
		} else if st, ok := status.FromError(uc.Error); ok && st.Code() < 17 {
			code = st.Code().String()
		}

		c.callDuration.With(prometheus.Labels{"code": code}).Observe(uc.Elapsed().Seconds())
		c.callCounter.With(prometheus.Labels{"code": code}).Inc()
	}
}

// collect memory usage by server workers
func (c *metricCollector) collectMemory(service roadrunner.Controllable, tick time.Duration) {
	started := false
	for {
		server := service.Server()
		if server == nil && started {
			// stopped
			return
		}

		started = true

		if workers, err := util.ServerState(server); err == nil {
			sum := 0.0
			for _, w := range workers {
				sum += float64(w.MemoryUsage)
			}

			c.workersMemory.Set(sum)
		}

		time.Sleep(tick)
	}
}
