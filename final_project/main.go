package main

import (
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/net"
)

type hostCollector struct {
    cpuMetric      *prometheus.Desc
    memMetric      *prometheus.Desc
    memTotalMetric *prometheus.Desc
    diskMetric     *prometheus.Desc
    diskFreeMetric *prometheus.Desc 
    netMetric      *prometheus.Desc
    netSentMetric  *prometheus.Desc 
}

func newHostCollector() *hostCollector {
    return &hostCollector{
        cpuMetric:  prometheus.NewDesc("host_cpu_usage_percent", "CPU usage percent of the VM", nil, nil),
        memMetric:  prometheus.NewDesc("host_memory_usage_bytes", "Memory used on the VM in bytes", nil, nil),
        memTotalMetric: prometheus.NewDesc("host_node_memory_MemTotal_bytes", "Total memory on the host node in bytes", nil, nil), 
        diskMetric: prometheus.NewDesc("host_disk_usage_percent", "Disk usage percent of root", nil, nil),
        
        diskFreeMetric: prometheus.NewDesc("host_node_disk_free_bytes", "Free disk space on root in bytes", nil, nil),
        netSentMetric:  prometheus.NewDesc("host_node_net_bytes_sent_total", "Total bytes sent by VM", nil, nil),

        netMetric:  prometheus.NewDesc("host_net_bytes_recv_total", "Total bytes received by VM", nil, nil),
    }
}

func (c *hostCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.cpuMetric
    ch <- c.memMetric
    ch <- c.memTotalMetric
    ch <- c.diskMetric
    ch <- c.diskFreeMetric 
    ch <- c.netMetric
    ch <- c.netSentMetric  
}

func (c *hostCollector) Collect(ch chan<- prometheus.Metric) {

    cpuPerc, _ := cpu.Percent(0, false)
    ch <- prometheus.MustNewConstMetric(c.cpuMetric, prometheus.GaugeValue, cpuPerc[0])

    v, _ := mem.VirtualMemory()
    ch <- prometheus.MustNewConstMetric(c.memMetric, prometheus.GaugeValue, float64(v.Used))
    ch <- prometheus.MustNewConstMetric(c.memTotalMetric, prometheus.GaugeValue, float64(v.Total)) 

    d, _ := disk.Usage("/")
    ch <- prometheus.MustNewConstMetric(c.diskMetric, prometheus.GaugeValue, d.UsedPercent)
    
    ch <- prometheus.MustNewConstMetric(c.diskFreeMetric, prometheus.GaugeValue, float64(d.Free))

    n, _ := net.IOCounters(false)
    if len(n) > 0 {
        ch <- prometheus.MustNewConstMetric(c.netMetric, prometheus.CounterValue, float64(n[0].BytesRecv))

        ch <- prometheus.MustNewConstMetric(c.netSentMetric, prometheus.CounterValue, float64(n[0].BytesSent))
    }
}

func main() {
    reg := prometheus.NewRegistry()

    m := newHostCollector()
    reg.MustRegister(m)

    http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

    log.Println("Exporter is running on :2112/metrics")
    log.Fatal(http.ListenAndServe(":2112", nil))
}