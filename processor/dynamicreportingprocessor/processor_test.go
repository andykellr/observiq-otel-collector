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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func testMetrics(t *testing.T, timestamp time.Time, offset time.Duration, values map[string]float64) pmetric.Metrics {
	md := pmetric.NewMetrics()
	metrics := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics()
	for name, value := range values {
		metric := metrics.AppendEmpty()
		metric.SetName(name)
		metric.SetEmptyGauge()
		dataPoint := metric.Gauge().DataPoints().AppendEmpty()
		dataPoint.SetDoubleValue(value)
		dataPoint.SetTimestamp(pcommon.Timestamp(timestamp.Add(offset).UnixNano()))
	}
	return md
}

func Test_sortIntervals(t *testing.T) {
	testCases := []struct {
		desc     string
		input    []Interval
		expected []Interval
	}{
		{
			desc: "Sort intervals",
			input: []Interval{
				{
					Threshold: 5.0,
					Interval:  1 * time.Hour,
				},
				{
					Threshold: 10.0,
					Interval:  1 * time.Second,
				},
				{
					Threshold: 0.0,
					Interval:  1 * time.Minute,
				},
				{
					Threshold: 15.0,
					Interval:  1 * time.Millisecond,
				},
			},
			expected: []Interval{
				{
					Threshold: 15.0,
					Interval:  1 * time.Millisecond,
				},
				{
					Threshold: 10.0,
					Interval:  1 * time.Second,
				},
				{
					Threshold: 5.0,
					Interval:  1 * time.Hour,
				},
				{
					Threshold: 0.0,
					Interval:  1 * time.Minute,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			p := newDynamicReportingProcessor(zap.NewNop(), &Config{
				Intervals: tc.input,
			})
			require.Equal(t, tc.expected, p.intervals)
		})
	}
}

func Test_processMetrics(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	init := testMetrics(t, now, 0, map[string]float64{
		"handle_metric": 0.0,
		"ignore_metric": 0.0,
	})

	intervals := []Interval{
		{
			Threshold: 10.0,
			Interval:  1 * time.Second,
		},
		{
			Threshold: 0.0,
			Interval:  1 * time.Minute,
		},
	}

	testCases := []struct {
		desc      string
		metrics   []string
		intervals []Interval
		init      pmetric.Metrics
		input     pmetric.Metrics
		expected  pmetric.Metrics
	}{
		{
			desc:      "Too soon, ignore",
			metrics:   []string{"handle_metric"},
			intervals: intervals,
			init:      init,
			input: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 5.0,
				"handle_metric": 5.0,
			}),
			expected: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 5.0,
			}),
		},
		{
			desc:      "Long time, include",
			metrics:   []string{"handle_metric"},
			intervals: intervals,
			init:      init,
			input: testMetrics(t, now, 1*time.Hour, map[string]float64{
				"ignore_metric": 5.0,
				"handle_metric": 5.0,
			}),
			expected: testMetrics(t, now, 1*time.Hour, map[string]float64{
				"ignore_metric": 5.0,
				"handle_metric": 5.0,
			}),
		},
		{
			desc:      "Spike, include",
			metrics:   []string{"handle_metric"},
			intervals: intervals,
			init:      init,
			input: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 50.0,
				"handle_metric": 50.0,
			}),
			expected: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 50.0,
				"handle_metric": 50.0,
			}),
		},
		{
			desc:      "Out of intervals, include",
			metrics:   []string{"handle_metric"},
			intervals: intervals,
			init:      init,
			input: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 5000.0,
				"handle_metric": 5000.0,
			}),
			expected: testMetrics(t, now, 10*time.Second, map[string]float64{
				"ignore_metric": 5000.0,
				"handle_metric": 5000.0,
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := &Config{
				Metrics:   tc.metrics,
				Intervals: tc.intervals,
			}

			processor := newDynamicReportingProcessor(zap.NewNop(), cfg)
			_, err := processor.processMetrics(context.Background(), tc.init)
			require.NoError(t, err)
			actual, err := processor.processMetrics(context.Background(), tc.input)
			require.NoError(t, err)
			require.ElementsMatch(t, asMetricsSlice(tc.expected), asMetricsSlice(actual))
		})
	}
}

func asMetricsSlice(metrics pmetric.Metrics) []pmetric.Metric {
	var names []pmetric.Metric
	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		for j := 0; j < metrics.ResourceMetrics().At(i).ScopeMetrics().Len(); j++ {
			for k := 0; k < metrics.ResourceMetrics().At(i).ScopeMetrics().At(j).Metrics().Len(); k++ {
				names = append(names, metrics.ResourceMetrics().At(i).ScopeMetrics().At(j).Metrics().At(k))
			}
		}
	}
	return names
}
