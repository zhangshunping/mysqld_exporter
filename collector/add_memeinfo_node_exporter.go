package collector

import (
	"context"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/prometheus/procfs"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"io"
	"strconv"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)


var (
	reParens = regexp.MustCompile(`\((.*)\)`)
	// The path of the proc filesystem.
	procPath   = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	sysPath    = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	rootfsPath = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
	)


type MeminfoCollector struct {
}

// Name of the Scraper. Should be unique.
func (MeminfoCollector) Name() string {
	return "node_exporter_meminfo"
}

// Help describes the role of the Scraper.
func (MeminfoCollector) Help() string {
	return "Collect the Node mem of all registered nodes"
}

// Version of MySQL from which scraper is available.
func (MeminfoCollector) Version() float64 {
	return 1.0
}




// Update calls (*meminfoCollector).getMemInfo to get the platform specific
// memory metrics.
func (c MeminfoCollector) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	var metricType prometheus.ValueType
	memInfo, err := c.getMemInfo()
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %w", err)
	}
	level.Debug(logger).Log("msg", "Set node_mem", "memInfo", memInfo)
	for k, v := range memInfo {
		if strings.HasSuffix(k, "_total") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("node", "memory", k),
				fmt.Sprintf("Memory information field %s.", k),
				nil, nil,
			),
			metricType, v,
		)
	}
	return nil
}



func (c *MeminfoCollector) getMemInfo() (map[string]float64, error) {
	file, err := os.Open(procFilePath("meminfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMemInfo(file)
}


func procFilePath(name string) string {
	return filepath.Join(*procPath, name)
}


func parseMemInfo(r io.Reader) (map[string]float64, error) {
	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		// Workaround for empty lines occasionally occur in CentOS 6.2 kernel 3.10.90.
		if len(parts) == 0 {
			continue
		}
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %w", err)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = reParens.ReplaceAllString(key, "_${1}")
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
			key = key + "_bytes"
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		memInfo[key] = fv
	}

	return memInfo, scanner.Err()
}

// check interface
var _ Scraper = MeminfoCollector{}

