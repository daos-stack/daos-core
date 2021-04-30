//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package promexp

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
)

var (
	_ prometheus.Metric    = &daosHistogram{}
	_ prometheus.Collector = &daosHistogramVec{}
)

type daosHistogram struct {
	desc      *prometheus.Desc
	name      string
	help      string
	sum       float64
	count     uint64
	labelVals []string
	buckets   map[float64]uint64
}

func newDaosHistogram(name, help string, desc *prometheus.Desc, labelVals []string) *daosHistogram {
	return &daosHistogram{
		name:      name,
		help:      help,
		desc:      desc,
		labelVals: labelVals,
		buckets:   make(map[float64]uint64),
	}
}

func (dh *daosHistogram) Desc() *prometheus.Desc {
	return dh.desc
}

func (dh *daosHistogram) sumBuckets() map[float64]uint64 {
	buckets := make([]float64, 0, len(dh.buckets))
	for key := range dh.buckets {
		buckets = append(buckets, key)
	}
	sort.Float64s(buckets)

	summedBuckets := make(map[float64]uint64)
	sum := uint64(0)
	for _, key := range buckets {
		sum += dh.buckets[key]
		summedBuckets[key] = sum
	}

	return summedBuckets
}

func (dh *daosHistogram) Write(out *dto.Metric) error {
	ch, err := prometheus.NewConstHistogram(dh.desc, dh.count, dh.sum, dh.sumBuckets(), dh.labelVals...)
	if err != nil {
		return err
	}

	return ch.Write(out)
}

func (dh *daosHistogram) AddBucketValue(bucket, value, sum float64, samples uint64) {
	dh.buckets[bucket] = uint64(value)
	//fmt.Printf("%f -> %f\n", value, bucket)
	dh.sum += sum
	dh.count += samples
}

type hashedMetricValue struct {
	metric    prometheus.Metric
	labelVals []string
}

type hashedMetrics map[uint64][]*hashedMetricValue

// daosHistogramVec is a simplified custom implementation of prometheus.HistogramVec.
// It is not designed for concurrency or currying.
type daosHistogramVec struct {
	opts       prometheus.HistogramOpts
	desc       *prometheus.Desc
	labelKeys  []string // stored here because prometheus.Desc is opaque to us
	histograms hashedMetrics
}

func (dhv *daosHistogramVec) Describe(ch chan<- *prometheus.Desc) {
	ch <- dhv.desc
}

func (dhv *daosHistogramVec) Collect(ch chan<- prometheus.Metric) {
	for _, histograms := range dhv.histograms {
		for _, histogram := range histograms {
			ch <- histogram.metric
		}
	}
}

func labelVals(labels prometheus.Labels) []string {
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	vals := make([]string, 0, len(labels))
	for _, key := range keys {
		vals = append(vals, labels[key])
	}
	return vals
}

func cmpLabelVals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (dhv *daosHistogramVec) addWithLabelValues(hashKey uint64, lvs []string) *hashedMetricValue {
	dh := newDaosHistogram(dhv.opts.Name, dhv.opts.Help, dhv.desc, lvs)
	hmv := &hashedMetricValue{
		metric:    dh,
		labelVals: lvs,
	}
	dhv.histograms[hashKey] = append(dhv.histograms[hashKey], hmv)

	return hmv
}

func (dhv *daosHistogramVec) GetWith(labels prometheus.Labels) (*daosHistogram, error) {
	hashKey := model.LabelsToSignature(labels)

	var hmv *hashedMetricValue
	lvs := labelVals(labels)
	_, found := dhv.histograms[hashKey]
	if !found {
		hmv = dhv.addWithLabelValues(hashKey, lvs)
	}

	if hmv == nil {
		for _, h := range dhv.histograms[hashKey] {
			if cmpLabelVals(h.labelVals, lvs) {
				hmv = h
				break
			}
		}

		if hmv == nil {
			hmv = dhv.addWithLabelValues(hashKey, lvs)
		}
	}

	dh, ok := hmv.metric.(*daosHistogram)
	if !ok {
		return nil, errors.New("stored something other than *daosHistogram")
	}
	return dh, nil
}

func (dhv *daosHistogramVec) With(labels prometheus.Labels) *daosHistogram {
	dh, err := dhv.GetWith(labels)
	if err != nil {
		panic(err)
	}
	return dh
}

func newDaosHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *daosHistogramVec {
	return &daosHistogramVec{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
			opts.Help,
			labelNames,
			opts.ConstLabels,
		),
		opts:       opts,
		labelKeys:  labelNames,
		histograms: make(hashedMetrics),
	}
}
