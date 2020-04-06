package metrics

import (
	"breeze/sensor"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Collector struct {
	temperature   *prometheus.GaugeVec
	thermalSensor sensor.Thermal
	nodeLabel     string
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("temperature", "thermal sensor metrics", nil, prometheus.Labels{
		"node_name": c.nodeLabel,
	})
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	read, err := c.thermalSensor.Read()
	if err != nil {
		log.Warnf("collector: unable to read temperature: %s", err.Error())
		return
	}
	gauge := c.temperature.WithLabelValues
	gauge(c.nodeLabel).Set(read)
	ch <- gauge(c.nodeLabel)
}

func New(nodeName string, thermalSensor sensor.Thermal) *Collector {
	return &Collector{
		nodeLabel:     nodeName,
		thermalSensor: thermalSensor,
		temperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "temperature",
		}, []string{"node"}),
	}
}

func (c *Collector) Run(addr string, port int) error {
	if port == 0 {
		port = 9999
	}

	if addr == "" {
		addr = "0.0.0.0"
	}

	if err := prometheus.Register(c); err != nil {
		return err
	}

	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
}
