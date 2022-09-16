package snowflake

import (
	"errors"
	"sync"
	"testing"
	"time"

	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

func TestInitGenerator(t *testing.T) {
	tests := []struct {
		name        string
		config      func() Config
		nodeID      uint64
		expectedErr error
	}{
		{
			name: "Success",
			config: func() Config {
				c, _ := InitConfig(defaultEpoch, defaultTimestampBits, 10)
				return c
			},
			nodeID:      12,
			expectedErr: nil,
		},
		{
			name: "Success - node-id with maximum allowed bit size",
			config: func() Config {
				c, _ := InitConfig(defaultEpoch, defaultTimestampBits, 3)
				return c
			},
			nodeID:      7,
			expectedErr: nil,
		},
		{
			name: "Invalid config",
			config: func() Config {
				return Config{}
			},
			nodeID:      12,
			expectedErr: errors.New("invalid config"),
		},
		{
			name: "Invalid node-id - larger than bit size",
			config: func() Config {
				c, _ := InitConfig(defaultEpoch, defaultTimestampBits, 3)
				return c
			},
			nodeID:      12,
			expectedErr: errors.New("nodeid can not be greater than [7] as per config"),
		},
		{
			name: "Invalid node-id - just larger than bit size",
			config: func() Config {
				c, _ := InitConfig(defaultEpoch, defaultTimestampBits, 3)
				return c
			},
			nodeID:      8,
			expectedErr: errors.New("nodeid can not be greater than [7] as per config"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitGenerator(tc.config(), tc.nodeID)
			if (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Error creating generator - actual - %v expected = %v", err, tc.expectedErr)
			}
		})
	}
}

type IDLogger struct {
	mu   sync.Mutex
	data map[uint64]bool
}

func (l *IDLogger) log(t *testing.T, id uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, exists := l.data[id]
	if exists {
		t.Errorf("duplicate ID - %d", id)
	}
	l.data[id] = true
	return nil
}

func TestGet(t *testing.T) {
	defaultConfig, _ := InitDefaultConfig()
	customeEpoch := uint64(time.Date(2015, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli())
	customConfig, _ := InitConfig(customeEpoch, 38, 14)
	tests := []struct {
		name         string
		getGenerator func() uid.Generator
		clientCount  int
		callCount    int
		expectedErr  error
	}{
		{
			name: "Success default config - single call by single client",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1,
			callCount:   1,
			expectedErr: nil,
		},
		{
			name: "Success default config - multiple calls by single client",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success default config - multiple calls by multiple clients",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1e2,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success custom config - single call by single client",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(customConfig, 4578)
				return gen
			},
			clientCount: 1,
			callCount:   1,
			expectedErr: nil,
		},
		{
			name: "Success custom config - multiple calls by single client",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(customConfig, 4578)
				return gen
			},
			clientCount: 1,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success custom config - multiple calls by multiple clients",
			getGenerator: func() uid.Generator {
				gen, _ := InitGenerator(customConfig, 4578)
				return gen
			},
			clientCount: 1e2,
			callCount:   1e2,
			expectedErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen := tc.getGenerator()
			logger := IDLogger{data: map[uint64]bool{}}
			var wg sync.WaitGroup
			wg.Add(tc.clientCount)
			for i := 0; i < tc.clientCount; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < tc.callCount; j++ {
						uid, err := gen.Get()
						logger.log(t, uid)
						if (tc.expectedErr != nil && err == nil) || (tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
							t.Errorf("Error calling generator get - actual - %v expected = %v", err, tc.expectedErr)
						}
					}
				}()
			}
			wg.Wait()
		})
	}
}

func BenchmarkGet(b *testing.B) {
	defaultConfig, _ := InitDefaultConfig()
	cases := []struct {
		name   string
		config func() Config
	}{
		{
			name:   "Default config - Tuned for high distribution of nodes (11 bits) and moderate throughput (10 bits) per node",
			config: func() Config { return defaultConfig },
		},
		{
			name: "Custom config - Tuned for low distribution of nodes (7 bits) and high throughput (14 bits) per node",
			config: func() Config {
				config, _ := InitConfig(defaultEpoch, 42, 7)
				return config
			},
		},
	}
	for _, bc := range cases {
		b.Run(bc.name, func(b *testing.B) {
			gen, _ := InitGenerator(bc.config(), 1)
			for i := 0; i < b.N; i++ {
				_, err := gen.Get()
				if err != nil {
					b.Errorf("Error occurred in benchmarking - %v", err)
				}
			}
		})
	}
}
