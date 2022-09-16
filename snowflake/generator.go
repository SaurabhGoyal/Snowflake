package snowflake

import (
	"fmt"
	"sync"
	"time"

	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

/*
Generator is a Snowflake based implmentation of UniqueId generation.
 - It calculates values in each node computationally.
 - It keeps only two values in its state - 1) last processed timestamp 2) last generated sequence id.
 - Each instance of generator within one application context must be assigned a node id.
 - Range of node ids will be dictated by the node-id bits given in the config.
*/
type Generator struct {
	mu     sync.Mutex
	nodeID uint64
	seqID  uint64
	lastTS uint64
	config Config
}

func (gen *Generator) getTimeStamp() uint64 {
	return uint64(time.Now().UnixMilli() - int64(gen.config.epoch))
}

/*
Get returns a unique id for current millisecond. If it can not provide a unique id, it waits till next millisecond and resets its sequence to provide a unique id.
*/
func (gen *Generator) Get() (uint64, error) {
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
func InitGenerator(config Config, nodeID uint64) (uid.Generator, error) {
	if config.timeStampBits <= 0 {
		return &Generator{}, fmt.Errorf("invalid config")
	}
	maxNodeID := uint64(1<<config.nodeIDBits - 1)
	if nodeID > maxNodeID {
		return &Generator{}, fmt.Errorf("nodeid can not be greater than [%d] as per config", maxNodeID)
	}
	return &Generator{
		nodeID: nodeID,
		seqID:  0,
		lastTS: config.epoch,
		config: config,
	}, nil
}
