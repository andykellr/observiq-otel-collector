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

package collector

import (
	"fmt"
	"os"

	"github.com/observiq/observiq-otel-collector/factories"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

const buildDescription = "observIQ's opentelemetry-collector distribution"

// NewSettings returns new settings for the collector with default values.
func NewSettings(configPaths []string, version string, loggingOpts []zap.Option) (*service.CollectorSettings, error) {
	factories, err := factories.DefaultFactories()
	if err != nil {
		return nil, fmt.Errorf("error while setting up default factories: %w", err)
	}

	buildInfo := component.BuildInfo{
		Command:     os.Args[0],
		Description: buildDescription,
		Version:     version,
	}

	fmp := fileprovider.New()
	configProviderSettings := service.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       configPaths,
			Providers:  map[string]confmap.Provider{fmp.Scheme(): fmp},
			Converters: []confmap.Converter{expandconverter.New()},
		},
	}
	provider, err := service.NewConfigProvider(configProviderSettings)
	if err != nil {
		return nil, err
	}

	return &service.CollectorSettings{
		Factories:               factories,
		BuildInfo:               buildInfo,
		LoggingOptions:          loggingOpts,
		ConfigProvider:          provider,
		DisableGracefulShutdown: true,
	}, nil
}
