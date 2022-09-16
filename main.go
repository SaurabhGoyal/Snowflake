package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	snowflake "github.com/SaurabhGoyal/Snowflake/snowflake"
	uid "github.com/SaurabhGoyal/Snowflake/uid"
)

const ITERS = int(1e2)
const CONSUMER_COUNT = int(3)
const SERVER_COUNT = byte(20)

type IDLogger struct {
	mu       sync.Mutex
	LastID   uint64
	IDS      map[uint64]bool
	OOOCount uint64
}

func (l *IDLogger) log(id uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.LastID > id {
		l.OOOCount += 1
	}
	_, exists := l.IDS[id]
	if exists {
		panic(fmt.Sprintf("conflict - %d - %v", id, l.IDS))
	}
	l.IDS[id] = true
	l.LastID = id
	// log.Printf("ID logged - [%d]", l.LastID)
	return nil
}

func consumeIds(uidGen uid.UIDGenerator, logger *IDLogger) {
	for i := 0; i < ITERS; i++ {
		uid, err := uidGen.Get()
		if err != nil {
			panic(err)
		}
		err = logger.log(uid)
		if err != nil {
			panic(err)
		}
		// time.Sleep(time.Millisecond * 200)
	}
}

func main() {
	log.Printf("Initialised main")
	config, err := snowflake.InitDefaultGeneratorConfig()
	log.Printf("Snowflake config - %v - err - %v", config, err)
	uidGen, _ := snowflake.InitGenerator(config, 2)
	idLogger := IDLogger{LastID: 0, IDS: map[uint64]bool{}}
	st := time.Now()
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		log.Printf("All done with OOO count - [%d] - Time taken - %d", idLogger.OOOCount, time.Since(st).Milliseconds())
	}()
	for i := 0; i < CONSUMER_COUNT; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumeIds(uidGen, &idLogger)
		}()
	}
}
