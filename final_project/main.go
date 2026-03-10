package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type hostCollector struct {
	cpuMetric  *prometheus.Desc
	memMetric  *prometheus.Desc
	diskMetric *prometheus.Desc
	netMetric  *prometheus.Desc
}

func newHostCollector() *hostCollector {
	return &hostCollector{
		cpuMetric:  prometheus.NewDesc("host_cpu_usage_percent", "CPU usage percent of the VM", nil, nil),
		memMetric:  prometheus.NewDesc("host_memory_usage_bytes", "Memory used on the VM in bytes", nil, nil),
		diskMetric: prometheus.NewDesc("host_disk_usage_percent", "Disk usage percent of root", nil, nil),
		netMetric:  prometheus.NewDesc("host_net_bytes_recv_total", "Total bytes received by VM", nil, nil),
	}
}

func (c *hostCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuMetric
	ch <- c.memMetric
	ch <- c.diskMetric
	ch <- c.netMetric
}

func (c *hostCollector) Collect(ch chan<- prometheus.Metric) {
	cpuPerc, _ := cpu.Percent(0, false)
	if len(cpuPerc) > 0 {
		ch <- prometheus.MustNewConstMetric(c.cpuMetric, prometheus.GaugeValue, cpuPerc[0])
	}

	v, _ := mem.VirtualMemory()
	if v != nil {
		ch <- prometheus.MustNewConstMetric(c.memMetric, prometheus.GaugeValue, float64(v.Used))
	}

	rootPath := os.Getenv("HOST_ROOTFS")
	if rootPath == "" {
		rootPath = "/"
	}
	d, _ := disk.Usage(rootPath)
	if d != nil {
		ch <- prometheus.MustNewConstMetric(c.diskMetric, prometheus.GaugeValue, d.UsedPercent)
	}

	n, _ := net.IOCounters(false)
	if len(n) > 0 {
		ch <- prometheus.MustNewConstMetric(c.netMetric, prometheus.CounterValue, float64(n[0].BytesRecv))
	}
}

func main() {
	reg := prometheus.NewRegistry()
	m := newHostCollector()
	reg.MustRegister(m)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}