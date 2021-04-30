//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package promexp

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/daos-stack/daos/src/control/lib/telemetry"
	"github.com/daos-stack/daos/src/control/logging"
)

type (
	bucketMap map[string]float64

	Collector struct {
		log            logging.Logger
		summary        *prometheus.SummaryVec
		ignoredMetrics []*regexp.Regexp
		sources        []*EngineSource
		buckets        bucketMap
	}

	CollectorOpts struct {
		Ignores   []string
		BucketMap bucketMap
	}

	EngineSource struct {
		ctx   context.Context
		Index uint32
		Rank  uint32
	}

	labelMap map[string]string
)

func (m bucketMap) Get(name string) (float64, error) {
	upper, found := m[name]
	if !found {
		return 0, errors.Errorf("%q not found in bucket map", name)
	}
	return upper, nil
}

func NewEngineSource(parent context.Context, idx uint32, rank uint32) (*EngineSource, error) {
	ctx, err := telemetry.Init(parent, idx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init telemetry")
	}

	return &EngineSource{
		ctx:   ctx,
		Index: idx,
		Rank:  rank,
	}, nil
}

func defaultBucketMap() bucketMap {
	return bucketMap{
		"256B":  1 << 8,
		"512B":  1 << 9,
		"1KB":   1 << 10,
		"2KB":   1 << 11,
		"4KB":   1 << 12,
		"8KB":   1 << 13,
		"16KB":  1 << 14,
		"32KB":  1 << 15,
		"64KB":  1 << 16,
		"128KB": 1 << 17,
		"256KB": 1 << 18,
		"512KB": 1 << 19,
		"1MB":   1 << 20,
		"2MB":   (1 << 20) * 2,
		"4MB":   (1 << 20) * 4,
		"_4MB":  math.Inf(0),
	}
}

func defaultCollectorOpts() *CollectorOpts {
	return &CollectorOpts{}
}

func NewCollector(log logging.Logger, opts *CollectorOpts, sources ...*EngineSource) (*Collector, error) {
	if len(sources) == 0 {
		return nil, errors.New("Collector must have > 0 sources")
	}

	if opts == nil {
		opts = defaultCollectorOpts()
	}
	if opts.BucketMap == nil {
		opts.BucketMap = defaultBucketMap()
	}

	c := &Collector{
		log:     log,
		sources: sources,
		buckets: opts.BucketMap,
		summary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "engine",
				Subsystem: "exporter",
				Name:      "scrape_duration_seconds",
				Help:      "daos_exporter: Duration of a scrape job.",
			},
			[]string{"source", "result"},
		),
	}

	for _, pat := range opts.Ignores {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compile %q", pat)
		}
		c.ignoredMetrics = append(c.ignoredMetrics, re)
	}

	return c, nil
}

func sanitizeMetricName(in string) string {
	return strings.Map(func(r rune) rune {
		switch {
		// Valid names for Prometheus are limited to:
		case r >= 'a' && r <= 'z': // lowercase letters
		case r >= 'A' && r <= 'Z': // uppercase letters
		case unicode.IsDigit(r): // digits
		default: // sanitize any other character
			return '_'
		}

		return r
	}, strings.TrimLeft(in, "/"))
}

func (es *EngineSource) Collect(log logging.Logger, ch chan<- *rankMetric) {
	metrics := make(chan telemetry.Metric)
	go func() {
		if err := telemetry.CollectMetrics(es.ctx, "", metrics); err != nil {
			log.Errorf("failed to collect metrics for engine rank %d: %s", es.Rank, err)
			return
		}
	}()

	for metric := range metrics {
		ch <- &rankMetric{
			r: es.Rank,
			m: metric,
		}
	}
}

type rankMetric struct {
	r uint32
	m telemetry.Metric
}

func (c *Collector) isIgnored(name string) bool {
	for _, re := range c.ignoredMetrics {
		if re.MatchString(name) {
			return true
		}
	}

	return false
}

func (lm labelMap) keys() (keys []string) {
	for label := range lm {
		keys = append(keys, label)
	}

	return
}

type gvMap map[string]*prometheus.GaugeVec

func (m gvMap) add(name, help string, value float64, labels labelMap) {
	var gv *prometheus.GaugeVec
	var found bool

	gv, found = m[name]
	if !found {
		gv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels.keys())
		m[name] = gv
	}
	gv.With(prometheus.Labels(labels)).Set(value)
}

type cvMap map[string]*prometheus.CounterVec

func (m cvMap) add(name, help string, value float64, labels labelMap) {
	var cv *prometheus.CounterVec
	var found bool

	cv, found = m[name]
	if !found {
		cv = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels.keys())
		m[name] = cv
	}
	cv.With(prometheus.Labels(labels)).Add(value)
}

type hvMap map[string]*daosHistogramVec

func (m hvMap) add(name, help string, bucket, value, sum float64, samples uint64, labels labelMap) {
	var hv *daosHistogramVec
	var found bool

	hv, found = m[name]
	if !found {
		hv = newDaosHistogramVec(prometheus.HistogramOpts{
			Name: name,
			Help: help,
		}, labels.keys())
		m[name] = hv
	}
	hv.With(prometheus.Labels(labels)).AddBucketValue(bucket, value, sum, samples)
}

type metricStat struct {
	name      string
	desc      string
	value     float64
	isCounter bool
}

func getMetricStats(baseName, desc string, ms telemetry.StatsMetric) (stats []*metricStat) {
	if ms.SampleSize() == 0 {
		return
	}

	for name, s := range map[string]struct {
		fn        func() float64
		desc      string
		isCounter bool
	}{
		"min": {
			fn:   ms.FloatMin,
			desc: " (min value)",
		},
		"max": {
			fn:   ms.FloatMax,
			desc: " (max value)",
		},
		"mean": {
			fn:   ms.Mean,
			desc: " (mean)",
		},
		"stddev": {
			fn:   ms.StdDev,
			desc: " (std dev)",
		},
		"samples": {
			fn:        func() float64 { return float64(ms.SampleSize()) },
			desc:      " (samples)",
			isCounter: true,
		},
	} {
		stats = append(stats, &metricStat{
			name:      baseName + "_" + name,
			desc:      desc + s.desc,
			value:     s.fn(),
			isCounter: s.isCounter,
		})
	}

	return
}

var (
	id_re  = regexp.MustCompile(`ID_+(\d+)_?`)
	io_re  = regexp.MustCompile(`io_+(\d+)_?`)
	net_re = regexp.MustCompile(`net_+(\d+)_+(\d+)_?`)
)

func fixPath(in string) (labels labelMap, name string) {
	name = sanitizeMetricName(in)

	labels = make(labelMap)

	// Clean up metric names and parse out useful labels
	name = id_re.ReplaceAllString(name, "")

	io_matches := io_re.FindStringSubmatch(name)
	if len(io_matches) > 0 {
		labels["target"] = io_matches[1]
		replacement := "io"
		if strings.HasSuffix(io_matches[0], "_") {
			replacement += "_"
		}
		name = io_re.ReplaceAllString(name, replacement)
	}

	net_matches := net_re.FindStringSubmatch(name)
	if len(net_matches) > 0 {
		labels["rank"] = net_matches[1]
		labels["context"] = net_matches[2]

		replacement := "net"
		if strings.HasSuffix(net_matches[0], "_") {
			replacement += "_"
		}
		name = net_re.ReplaceAllString(name, replacement)
	}

	return
}

type patternSet []*regexp.Regexp

func (ps patternSet) Matches(in string) (matches []string) {
	for _, p := range ps {
		matches = p.FindStringSubmatch(in)
		if len(matches) > 0 {
			return
		}
	}
	return
}

var (
	latHistPats = patternSet{
		regexp.MustCompile(`(engine_io_fetch_latency)_(\w+)`),
		regexp.MustCompile(`(engine_io_update_latency)_(\w+)`),
	}
)

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	rankMetrics := make(chan *rankMetric)
	go func(sources []*EngineSource) {
		for _, source := range c.sources {
			source.Collect(c.log, rankMetrics)
		}
		close(rankMetrics)
	}(c.sources)

	gauges := make(gvMap)
	counters := make(cvMap)
	histograms := make(hvMap)

	for rm := range rankMetrics {
		labels, path := fixPath(rm.m.Path())
		labels["rank"] = fmt.Sprintf("%d", rm.r)

		name := sanitizeMetricName(rm.m.Name())

		baseName := prometheus.BuildFQName("engine", path, name)
		desc := rm.m.Desc()

		if matches := latHistPats.Matches(baseName); len(matches) > 0 {
			sm, ok := rm.m.(telemetry.StatsMetric)
			if !ok || sm.FloatValue() == 0 {
				continue
			}

			baseName = matches[1]
			upperBound, err := c.buckets.Get(matches[2])
			if err != nil {
				c.log.Errorf("unable to parse bucket name %q", matches[2])
				continue
			}

			// Create a histogram for per-size latencies.
			histograms.add(baseName, desc, upperBound, sm.FloatValue(), sm.FloatSum(), sm.SampleSize(), labels)

			// Create a histogram for size distribution.
			sizeName := strings.ReplaceAll(baseName, "latency", "sizes")
			histograms.add(sizeName, desc, upperBound, float64(sm.SampleSize()), sm.FloatSum(), sm.SampleSize(), labels)

			continue
		}

		switch metric := rm.m.(type) {
		case *telemetry.Gauge:
			if c.isIgnored(baseName) {
				continue
			}

			gauges.add(baseName, desc, rm.m.FloatValue(), labels)
			for _, ms := range getMetricStats(baseName, desc, metric) {
				if ms.isCounter {
					counters.add(ms.name, ms.desc, ms.value, labels)
					continue
				}
				gauges.add(ms.name, ms.desc, ms.value, labels)
			}
		case *telemetry.Counter:
			if c.isIgnored(baseName) {
				break
			}

			counters.add(baseName, desc, rm.m.FloatValue(), labels)
		default:
			c.log.Errorf("metric type %d not supported", rm.m.Type())
		}
	}

	for _, gv := range gauges {
		gv.Collect(ch)
	}
	for _, cv := range counters {
		cv.Collect(ch)
	}
	for _, hv := range histograms {
		hv.Collect(ch)
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.summary.Describe(ch)
}
