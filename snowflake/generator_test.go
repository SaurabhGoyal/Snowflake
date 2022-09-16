package snowflake

import (
	"errors"
	"sync"
	"testing"
	"time"

	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

func validateGen(t *testing.T, config generatorConfig, err error, expectedTimeStampShift uint64, expectedNodeIdShift uint64, expectedErr error) {
	if (expectedErr != nil && err == nil) || (expectedErr != nil && err.Error() != expectedErr.Error()) {
		t.Errorf("Error creating config - actual - %v expected = %v", err, expectedErr)
	}
	if err != nil {
		return
	}
	if config.timeStampShift != expectedTimeStampShift {
		t.Errorf("timestamp shift mismatch - actual - %d, expected - %d", config.timeStampShift, expectedTimeStampShift)
	}
	if config.nodeIdShift != expectedNodeIdShift {
		t.Errorf("nodeId shift mismatch - actual - %d, expected - %d", config.nodeIdShift, expectedNodeIdShift)
	}
}

func TestInitGenerator(t *testing.T) {
	tests := []struct {
		name        string
		config      generatorConfig
		nodeId      uint64
		expectedErr error
	}{
		{
			name: "Success",
			config: generatorConfig{
				nodeIdBits: 10,
			},
			nodeId:      12,
			expectedErr: nil,
		},
		{
			name: "Success - node-id with maximum allowed bit size",
			config: generatorConfig{
				nodeIdBits: 3,
			},
			nodeId:      7,
			expectedErr: nil,
		},
		{
			name: "Invalid node-id - larger than bit size",
			config: generatorConfig{
				nodeIdBits: 3,
			},
			nodeId:      12,
			expectedErr: errors.New("nodeid can not be greater than [7] as per config"),
		},
		{
			name: "Invalid node-id - just larger than bit size",
			config: generatorConfig{
				nodeIdBits: 3,
			},
			nodeId:      8,
			expectedErr: errors.New("nodeid can not be greater than [7] as per config"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitGenerator(tc.config, tc.nodeId)
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
	defaultConfig, _ := InitDefaultGeneratorConfig()
	customeEpoch := uint64(time.Date(2015, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli())
	customConfig, _ := InitGeneratorConfig(customeEpoch, 38, 14)
	tests := []struct {
		name         string
		getGenerator func() uid.UIDGenerator
		clientCount  int
		callCount    int
		expectedErr  error
	}{
		{
			name: "Success default config - single call by single client",
			getGenerator: func() uid.UIDGenerator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1,
			callCount:   1,
			expectedErr: nil,
		},
		{
			name: "Success default config - multiple calls by single client",
			getGenerator: func() uid.UIDGenerator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success default config - multiple calls by multiple clients",
			getGenerator: func() uid.UIDGenerator {
				gen, _ := InitGenerator(defaultConfig, 3)
				return gen
			},
			clientCount: 1e2,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success custom config - single call by single client",
			getGenerator: func() uid.UIDGenerator {
				gen, _ := InitGenerator(customConfig, 4578)
				return gen
			},
			clientCount: 1,
			callCount:   1,
			expectedErr: nil,
		},
		{
			name: "Success custom config - multiple calls by single client",
			getGenerator: func() uid.UIDGenerator {
				gen, _ := InitGenerator(customConfig, 4578)
				return gen
			},
			clientCount: 1,
			callCount:   1e2,
			expectedErr: nil,
		},
		{
			name: "Success custom config - multiple calls by multiple clients",
			getGenerator: func() uid.UIDGenerator {
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
