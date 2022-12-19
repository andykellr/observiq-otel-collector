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

// Package dynamicreportingprocessor provides a processor that dynamically adjusts the rate of metrics based on the value.
package dynamicreportingprocessor

import (
	"time"

	"go.opentelemetry.io/collector/config"
)

// Interval represents a single reporting interval
type Interval struct {
	Threshold float64       `mapstructure:"threshold"`
	Interval  time.Duration `mapstructure:"interval"`
}

// Config is the configuration for the processor
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	Metrics   []string   `mapstructure:"metrics"`
	Intervals []Interval `mapstructure:"intervals"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	// TODO(observIQ): validate the config
	return nil
}
