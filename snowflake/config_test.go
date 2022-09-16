package snowflake

import (
	"errors"
	"testing"
	"time"
)

func validateConfig(t *testing.T, config generatorConfig, err error, expectedTimeStampShift uint64, expectedNodeIdShift uint64, expectedErr error) {
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

func TestInitDefaultGeneratorConfig(t *testing.T) {
	tests := []struct {
		name                   string
		expectedTimeStampShift uint64
		expectedNodeIdShift    uint64
		expectedErr            error
	}{
		{
			name:                   "Default",
			expectedTimeStampShift: 19,
			expectedNodeIdShift:    9,
			expectedErr:            nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := InitDefaultGeneratorConfig()
			validateConfig(t, config, err, tc.expectedTimeStampShift, tc.expectedNodeIdShift, tc.expectedErr)
		})
	}
}

func TestInitGeneratorConfig(t *testing.T) {
	tests := []struct {
		name                   string
		epoch                  uint64
		timestampBits          uint64
		nodeIdBits             uint64
		expectedTimeStampShift uint64
		expectedNodeIdShift    uint64
		expectedErr            error
	}{
		{
			name:                   "Success",
			epoch:                  uint64(time.Date(2020, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIdBits:             6,
			expectedTimeStampShift: 12,
			expectedNodeIdShift:    6,
			expectedErr:            nil,
		},
		{
			name:                   "Invalid bits",
			epoch:                  uint64(time.Date(2020, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIdBits:             16,
			expectedTimeStampShift: 0,
			expectedNodeIdShift:    0,
			expectedErr:            errors.New("timestamp and nodeid can accommodate maximum [59] bits"),
		},
		{
			name:                   "Invalid epoch",
			epoch:                  uint64(time.Date(time.Now().Year()+2, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIdBits:             16,
			expectedTimeStampShift: 0,
			expectedNodeIdShift:    0,
			expectedErr:            errors.New("epoch must be in past - given epoch [1734047999000] is in future"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := InitGeneratorConfig(tc.epoch, tc.timestampBits, tc.nodeIdBits)
			validateConfig(t, config, err, tc.expectedTimeStampShift, tc.expectedNodeIdShift, tc.expectedErr)
		})
	}
}
