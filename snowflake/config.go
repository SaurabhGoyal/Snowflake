package snowflake

import (
	"fmt"
	"time"
)

// time.Date(2010, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()
const defaultEpoch = uint64(1292198399000)
const uidBitLength = uint64(63)
const minAllowedBitsForTimestamp = uint64(40)
const maxAllowedBitsForTimestampAndNodeID = uint64(59)
const defaultTimestampBits = uint64(42)
const defaultNodeIDBits = uint64(11)
const millisecondInYear = 365 * 86400000

/*
Config contains the user given parameters and required values derived from those to dictate how unique-id generation process.
*/
type Config struct {
	epoch          uint64
	timeStampBits  uint64
	timeStampShift uint64
	nodeIDBits     uint64
	nodeIDShift    uint64
}

/*
InitConfig generates a config to be used by generator to generate IDs.
 - Use latest possible epoch based on what timestamps are irrelevant by your app. This can increase the years that you can target.
 - Use higher or lower timestamp bits based on how long (generally in years) the lifecycle of a generated unique id should be.
 - Use higher or lower nodeID bits based on how many servers are going to be involved in unique id generation.
 - Above two values directly impact the range of unique ids that one server can generate per millisecond while being within the constraint of 64 bit.
Choose above wisely as higher range per server gives better performance in high throughput systems.
*/
func InitConfig(epoch uint64, timestampBits uint64, nodeIDBits uint64) (Config, error) {
	currentTs := uint64(time.Now().UnixMilli())
	if epoch >= currentTs {
		return Config{}, fmt.Errorf("epoch must be in past - given epoch [%d] is in future", epoch)
	}
	if timestampBits < minAllowedBitsForTimestamp {
		return Config{}, fmt.Errorf("timestamp length must be atleast [%d] bits to be able to generate unique ids", minAllowedBitsForTimestamp)
	}
	if timestampBits+nodeIDBits > maxAllowedBitsForTimestampAndNodeID {
		return Config{}, fmt.Errorf("timestamp and nodeid can accommodate maximum [%d] bits", maxAllowedBitsForTimestampAndNodeID)
	}
	return Config{
		epoch:          epoch,
		timeStampBits:  timestampBits,
		timeStampShift: uidBitLength - timestampBits,
		nodeIDBits:     nodeIDBits,
		nodeIDShift:    uidBitLength - (timestampBits + nodeIDBits),
	}, nil
}

/*
InitDefaultConfig generates a default config suited to current applications.
 - Epoch - 2010/12/12
 - Timestamp bits - 42 (~139 years of id uniqueness)
 - NodeID bits - 11 (2048 servers)
This in turn means 1048 numbers per milliosecond per server can be generated uniquely.
*/
func InitDefaultConfig() (Config, error) {
	return InitConfig(defaultEpoch, defaultTimestampBits, defaultNodeIDBits)
}

func (config *Config) String() string {
	maxTPS := 1 << (config.nodeIDBits + config.nodeIDShift)
	maxServers := 1 << config.nodeIDBits
	maxTPSPerServer := 1 << config.nodeIDShift
	maxLifeTimeMS := 1 << config.timeStampBits
	maxLifeTime := float64(maxLifeTimeMS) / millisecondInYear
	output :=
		`
==============================================================================
Initialising Snowflake Unique ID Generator Config
==============================================================================
Config (64 bit ID)
+------------------------------------------------------------------------+
| 1 Bit Unused | %d Bit Timestamp |  %d Bit NodeID  | %d Bit Sequence ID |
+------------------------------------------------------------------------+
Output (TP = requests per millisecond)
+---------------------------------------------------------------------------------+
| %0.2f Years of uniqueness lifetime | %d TP across %d servers  | %d TP per server |
+---------------------------------------------------------------------------------+
==============================================================================
`
	return fmt.Sprintf(output, config.timeStampBits, config.nodeIDBits, config.nodeIDShift, maxLifeTime, maxTPS, maxServers, maxTPSPerServer)
}
