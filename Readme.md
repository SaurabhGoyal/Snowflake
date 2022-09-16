# Snowflake
Snowflake is a go package that provides a simple implementation of unique-id generation based on [Twitter's snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) logic.
```
+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
+--------------------------------------------------------------------------+
```

## Status
Project has just started and is not usable in any environment as it is missing edge case testing and benchmarking.

## Logic
- There are two components - `generator` and `generatorConfig`. Both are private to prevent initialisation with unexpected values. Use exported constructors which perform validations.
- `generatorConfig` has a default implementation which provides 32768 unique ids per millisecond till year 2079 (~69 years from epoch). These values can be customised to support higher throughput or longer period of valid generation of unique IDs. Default config uses following details -
  - epoch - `2010/12/12/23/59/59/0 UTC`.
  - timestampBits - 48
  - nodeIdBits - 10

## Hands On
- Install (Assuming you already have a working Go environment, if not please see [this page](https://go.dev/doc/install) first)
```
go get github.com/SaurabhGoyal/Snowflake
```
- Use
```go
package main

import (
	"log"

snowflake "github.com/SaurabhGoyal/Snowflake/snowflake"
)

func main() {
	// With default config
    config, _ := snowflake.InitDefaultGeneratorConfig()
	uidGen, _ := snowflake.InitGenerator(config, 1)
	uid, err := uidGen.Get()
    if err != nil {
        panic(err)
    }
    log.Printf("UID - %d", uid)

    // With custom config
    epoch := uint64(time.Date(2015, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli())
    config, _ = snowflake.InitGeneratorConfig(epoch, 50, 7)
	uidGen, _ = snowflake.InitGenerator(config, 1)
	uid, err = uidGen.Get()
    if err != nil {
        panic(err)
    }
    log.Printf("UID - %d", uid)  
}

```