package snowflake

import (
	"errors"
	"fmt"
	"time"
)

// time.Date(2010, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli()
const DEFAULT_EPOCH = 1292198399000
const ID_BITS = uint64(63)
const MAX_ALLOWED_BITS_FOR_TIMESTAMP_AND_NODE_ID = uint64(59)
const DEFAULT_TIMESTAMP_BITS = uint64(48)
const DEFAULT_NODE_ID_BITS = uint64(10)

type generatorConfig struct {
	epoch          uint64
	timeStampBits  uint64
	timeStampShift uint64
	nodeIdBits     uint64
	nodeIdShift    uint64
}

func InitGeneratorConfig(epoch uint64, timestampBits uint64, nodeIdBits uint64) (generatorConfig, error) {
	current_ts := uint64(time.Now().UnixMilli())
	if epoch >= current_ts {
		return generatorConfig{}, errors.New(
			fmt.Sprintf("epoch must be in past - given epoch [%d] is in future", epoch),
		)
	}
	if timestampBits+nodeIdBits > MAX_ALLOWED_BITS_FOR_TIMESTAMP_AND_NODE_ID {
		return generatorConfig{}, errors.New(
			fmt.Sprintf("timestamp and nodeid can accommodate maximum [%d] bits", MAX_ALLOWED_BITS_FOR_TIMESTAMP_AND_NODE_ID),
		)
	}
	return generatorConfig{
		epoch:          epoch,
		timeStampBits:  timestampBits,
		timeStampShift: ID_BITS - timestampBits,
		nodeIdBits:     nodeIdBits,
		nodeIdShift:    ID_BITS - (timestampBits + nodeIdBits),
	}, nil
}

func InitDefaultGeneratorConfig() (generatorConfig, error) {
	return InitGeneratorConfig(uint64(DEFAULT_EPOCH), DEFAULT_TIMESTAMP_BITS, DEFAULT_NODE_ID_BITS)
}
