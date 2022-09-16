package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"

	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

type generator struct {
	mu     sync.Mutex
	nodeId uint64
	seqId  uint64
	lastTS uint64
	config generatorConfig
}

func (gen *generator) getTimeStamp() uint64 {
	return uint64(time.Now().UnixMilli() - int64(gen.config.epoch))
}

func (gen *generator) Get() (uint64, error) {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	seqId := uint64(0)
	timestamp := gen.getTimeStamp()
	if timestamp <= gen.lastTS {
		seqId = gen.seqId + 1
		if seqId > (1<<gen.config.nodeIdShift - 1) {
			time.Sleep(time.Millisecond)
			timestamp = gen.getTimeStamp()
			seqId = 0
		}
	}
	gen.lastTS = timestamp
	gen.seqId = seqId
	uid := timestamp<<gen.config.timeStampShift | gen.nodeId<<gen.config.nodeIdShift | seqId
	return uid, nil
}

func InitGenerator(config generatorConfig, nodeId uint64) (uid.UIDGenerator, error) {
	maxNodeId := uint64(1<<config.nodeIdBits - 1)
	if nodeId > maxNodeId {
		return &generator{}, errors.New(
			fmt.Sprintf("nodeid can not be greater than [%d] as per config", maxNodeId),
		)
	}
	return &generator{
		nodeId: nodeId,
		seqId:  0,
		lastTS: config.epoch,
		config: config,
	}, nil
}
