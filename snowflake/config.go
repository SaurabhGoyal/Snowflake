package snowflake

import (
	"fmt"
	"time"
)

// time.Date(2010, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()
const defaultEpoch = 1292198399000
const uidBitLength = uint64(63)
const minAllowedBitsForTimestamp = uint64(40)
const maxAllowedBitsForTimestampAndNodeID = uint64(59)
const defaultTimestampBits = uint64(42)
const defaultNodeIDBits = uint64(11)

type generatorConfig struct {
	epoch          uint64
	timeStampBits  uint64
	timeStampShift uint64
	nodeIDBits     uint64
	nodeIDShift    uint64
}

/*
InitGeneratorConfig generates a config to be used by generator to generate IDs.
 - Use latest possible epoch based on what timestamps are irrelevant by your app. This can increase the years that you can target.
 - Use higher or lower timestamp bits based on how long (generally in years) the lifecycle of a generated unique id should be.
 - Use higher or lower nodeID bits based on how many servers are going to be involved in unique id generation.
 - Above two values directly impact the range of unique ids that one server can generate per millisecond while being within the constraint of 64 bit.
Choose above wisely as higher range per server gives better performance in high throughput systems.
*/
func InitGeneratorConfig(epoch uint64, timestampBits uint64, nodeIDBits uint64) (generatorConfig, error) {
	currentTs := uint64(time.Now().UnixMilli())
	if epoch >= currentTs {
		return generatorConfig{}, fmt.Errorf("epoch must be in past - given epoch [%d] is in future", epoch)
	}
	if timestampBits < minAllowedBitsForTimestamp {
		return generatorConfig{}, fmt.Errorf("timestamp length must be atleast [%d] bits to be able to generate unique ids", minAllowedBitsForTimestamp)
	}
	if timestampBits+nodeIDBits > maxAllowedBitsForTimestampAndNodeID {
		return generatorConfig{}, fmt.Errorf("timestamp and nodeid can accommodate maximum [%d] bits", maxAllowedBitsForTimestampAndNodeID)
	}
	return generatorConfig{
		epoch:          epoch,
		timeStampBits:  timestampBits,
		timeStampShift: uidBitLength - timestampBits,
		nodeIDBits:     nodeIDBits,
		nodeIDShift:    uidBitLength - (timestampBits + nodeIDBits),
	}, nil
}

/*
InitDefaultGeneratorConfig generates a default config suited to current applications.
 - Epoch - 2010/12/12
 - Timestamp bits - 42 (~139 years of id uniqueness)
 - NodeID bits - 11 (2048 servers)
This in turn means 1048 numbers per milliosecond per server can be generated uniquely.
*/
func InitDefaultGeneratorConfig() (generatorConfig, error) {
	return InitGeneratorConfig(uint64(defaultEpoch), defaultTimestampBits, defaultNodeIDBits)
}
