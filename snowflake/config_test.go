package snowflake

import (
	"errors"
	"testing"
	"time"
)

func validateConfig(t *testing.T, config Config, err error, expectedTimeStampShift uint64, expectedNodeIDShift uint64, expectedErr error) {
	if (expectedErr != nil && err == nil) || (expectedErr != nil && err.Error() != expectedErr.Error()) {
		t.Errorf("Error creating config - actual - %v expected = %v", err, expectedErr)
	}
	if err != nil {
		return
	}
	if config.timeStampShift != expectedTimeStampShift {
		t.Errorf("timestamp shift mismatch - actual - %d, expected - %d", config.timeStampShift, expectedTimeStampShift)
	}
	if config.nodeIDShift != expectedNodeIDShift {
		t.Errorf("nodeId shift mismatch - actual - %d, expected - %d", config.nodeIDShift, expectedNodeIDShift)
	}
}

func TestInitDefaultGeneratorConfig(t *testing.T) {
	tests := []struct {
		name                   string
		expectedTimeStampShift uint64
		expectedNodeIDShift    uint64
		expectedErr            error
	}{
		{
			name:                   "Default",
			expectedTimeStampShift: 21,
			expectedNodeIDShift:    10,
			expectedErr:            nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := InitDefaultConfig()
			validateConfig(t, config, err, tc.expectedTimeStampShift, tc.expectedNodeIDShift, tc.expectedErr)
		})
	}
}

func TestInitGeneratorConfig(t *testing.T) {
	tests := []struct {
		name                   string
		epoch                  uint64
		timestampBits          uint64
		nodeIDBits             uint64
		expectedTimeStampShift uint64
		expectedNodeIDShift    uint64
		expectedErr            error
	}{
		{
			name:                   "Success",
			epoch:                  uint64(time.Date(2020, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIDBits:             6,
			expectedTimeStampShift: 12,
			expectedNodeIDShift:    6,
			expectedErr:            nil,
		},
		{
			name:                   "Invalid bits - Greater than limit",
			epoch:                  uint64(time.Date(2020, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIDBits:             11,
			expectedTimeStampShift: 0,
			expectedNodeIDShift:    0,
			expectedErr:            errors.New("timestamp and nodeid can accommodate maximum [59] bits"),
		},
		{
			name:                   "Invalid bits - Less than required for uniqueness",
			epoch:                  uint64(time.Date(2020, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          31,
			nodeIDBits:             11,
			expectedTimeStampShift: 0,
			expectedNodeIDShift:    0,
			expectedErr:            errors.New("timestamp length must be atleast [40] bits to be able to generate unique ids"),
		},
		{
			name:                   "Invalid epoch",
			epoch:                  uint64(time.Date(time.Now().Year()+2, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()),
			timestampBits:          51,
			nodeIDBits:             16,
			expectedTimeStampShift: 0,
			expectedNodeIDShift:    0,
			expectedErr:            errors.New("epoch must be in past - given epoch [1734047999000] is in future"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := InitConfig(tc.epoch, tc.timestampBits, tc.nodeIDBits)
			validateConfig(t, config, err, tc.expectedTimeStampShift, tc.expectedNodeIDShift, tc.expectedErr)
		})
	}
}
