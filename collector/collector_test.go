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
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectorRunValid(t *testing.T) {
	ctx := context.Background()

	collector := New([]string{"./test/valid.yaml"}, "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorRunMultiple(t *testing.T) {
	collector := New([]string{"./test/valid.yaml"}, "0.0.0", nil)
	for i := 1; i < 5; i++ {
		ctx := context.Background()

		attempt := fmt.Sprintf("Attempt %d", i)
		t.Run(attempt, func(t *testing.T) {
			err := collector.Run(ctx)
			require.NoError(t, err)

			status := <-collector.Status()
			require.True(t, status.Running)
			require.NoError(t, status.Err)

			collector.Stop()
			status = <-collector.Status()
			require.False(t, status.Running)
		})
	}
}

func TestCollectorRunInvalidConfig(t *testing.T) {
	ctx := context.Background()

	collector := New([]string{"./test/invalid.yaml"}, "0.0.0", nil)
	err := collector.Run(ctx)
	require.Error(t, err)

	status := <-collector.Status()
	require.False(t, status.Running)
	require.Error(t, status.Err)
	require.ErrorContains(t, status.Err, "cannot unmarshal the configuration")
}

// There currently exists a limitation in the collector lifecycle regarding context.
// Context is not respected when starting the collector and a collector could run indefinitely
// in this scenario. Once this is addressed, we can readd this test.
//
// func TestCollectorRunCancelledContext(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	cancel()

// 	collector := New("./test/valid.yaml", "0.0.0", nil)
// 	err := collector.Run(ctx)
// 	require.EqualError(t, context.Canceled, err.Error())
// }

func TestCollectorRunTwice(t *testing.T) {
	ctx := context.Background()

	collector := New([]string{"./test/valid.yaml"}, "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)
	defer collector.Stop()

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	err = collector.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service already running")

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorRestart(t *testing.T) {
	ctx := context.Background()

	collector := New([]string{"./test/valid.yaml"}, "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	err = collector.Restart(ctx)
	require.NoError(t, err)

	status = <-collector.Status()
	require.True(t, status.Running)

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorPrematureStop(t *testing.T) {
	collector := New([]string{"./test/valid.yaml"}, "0.0.0", nil)
	collector.Stop()
	require.Equal(t, 0, len(collector.Status()))
}
