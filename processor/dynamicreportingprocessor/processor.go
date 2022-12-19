// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamicreportingprocessor

import (
	"context"
	"sort"
	"time"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type metricState struct {
	lastReportedAt    time.Time
	lastIntervalIndex int
}

type dynamicReportingProcessor struct {
	logger    *zap.Logger
	intervals []Interval
	state     map[string]*metricState
}

func newDynamicReportingProcessor(logger *zap.Logger, cfg *Config) *dynamicReportingProcessor {
	state := map[string]*metricState{}
	for _, m := range cfg.Metrics {
		state[m] = &metricState{}
	}

	// sort by threshold in descending order
	intervals := cfg.Intervals
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Threshold > intervals[j].Threshold
	})

	return &dynamicReportingProcessor{
		logger:    logger,
		state:     state,
		intervals: intervals,
	}
}

func (d *dynamicReportingProcessor) intervalFromValue(value float64) (time.Duration, int) {
	for i, interval := range d.intervals {
		if value >= interval.Threshold {
			return interval.Interval, i
		}
	}
	return 0, -1
}

func (d *dynamicReportingProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		for j := 0; j < md.ResourceMetrics().At(i).ScopeMetrics().Len(); j++ {
			md.ResourceMetrics().At(i).ScopeMetrics().At(j).Metrics().RemoveIf(func(m pmetric.Metric) bool {
				// this only works for gauge metrics right now
				if m.Type() != pmetric.MetricTypeGauge {
					return false
				}
				if state, ok := d.state[m.Name()]; ok {
					dataPoint := m.Gauge().DataPoints().At(0)
					interval, intervalIndex := d.intervalFromValue(dataPoint.DoubleValue())
					metricTime := dataPoint.Timestamp().AsTime()

					// we only want to skip if we are in the same interval and the metric time is less than the interval
					skip := intervalIndex == state.lastIntervalIndex && metricTime.Sub(state.lastReportedAt) < interval

					// fmt.Printf("metric: %s, value: %f, interval: %s, metricTime: %s, lastReportedAt: %s, skip: %t\n",
					// 	m.Name(),
					// 	dataPoint.DoubleValue(),
					// 	interval,
					// 	metricTime,
					// 	state.lastReportedAt,
					// 	skip,
					// )

					if skip {
						return true
					}
					state.lastReportedAt = metricTime
					state.lastIntervalIndex = intervalIndex
				}
				return false
			})
		}
	}
	return md, nil
}
