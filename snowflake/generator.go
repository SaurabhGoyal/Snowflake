package snowflake

import (
	"fmt"
	"sync"
	"time"

	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

type generator struct {
	mu     sync.Mutex
	nodeID uint64
	seqID  uint64
	lastTS uint64
	config generatorConfig
}

func (gen *generator) getTimeStamp() uint64 {
	return uint64(time.Now().UnixMilli() - int64(gen.config.epoch))
}

func (gen *generator) Get() (uint64, error) {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	seqID := uint64(0)
	timestamp := gen.getTimeStamp()
	if timestamp <= gen.lastTS {
		seqID = gen.seqID + 1
		if seqID > ((1 << gen.config.nodeIDShift) - 1) {
			time.Sleep(time.Millisecond)
			timestamp = gen.getTimeStamp()
			seqID = 0
		}
	}
	gen.lastTS = timestamp
	gen.seqID = seqID
	uid := timestamp<<gen.config.timeStampShift | gen.nodeID<<gen.config.nodeIDShift | seqID
	return uid, nil
}

/*
InitGenerator initialises a unique id generator as per Snowflake algo using given config.
*/
func InitGenerator(config generatorConfig, nodeID uint64) (uid.Generator, error) {
	maxNodeID := uint64(1<<config.nodeIDBits - 1)
	if nodeID > maxNodeID {
		return &generator{}, fmt.Errorf("nodeid can not be greater than [%d] as per config", maxNodeID)
	}
	return &generator{
		nodeID: nodeID,
		seqID:  0,
		lastTS: config.epoch,
		config: config,
	}, nil
}
