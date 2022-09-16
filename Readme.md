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
- `generatorConfig` has a default implementation which provides 2 ^ 21 (2097152) unique ids per millisecond (1024 per node, across 2048 nodes) till year 2149 (~139 years from epoch). These values can be customised to support higher throughput or longer period of valid generation of unique IDs. Default config uses following details -
  - epoch - `2010/12/12/23/59/59/0 UTC`.
  - timestampBits - 42
  - nodeIdBits - 11

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

## Benchmarks
- Result for one machine
```
➜  uid_generator git:(main) ✗ go test -v -run=^# -bench=. -count=4 ./... -benchmem

?       github.com/SaurabhGoyal/Snowflake       [no test files]
goos: darwin
goarch: amd64
pkg: github.com/SaurabhGoyal/Snowflake/snowflake
cpu: VirtualApple @ 2.50GHz
BenchmarkGet
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8                1000000           1158 ns/op        0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8                 896391           1163 ns/op        0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8                1000000              1187 ns/op               0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8                1000000              1209 ns/op               0 B/op            0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                      11487432               103.1 ns/op             0 B/op            0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                      12016778               103.7 ns/op             0 B/op            0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                      11961061               102.9 ns/op             0 B/op            0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                      12601296               102.4 ns/op             0 B/op            0 allocs/op
PASS
ok      github.com/SaurabhGoyal/Snowflake/snowflake     11.124s
?       github.com/SaurabhGoyal/Snowflake/uid   [no test files]
➜  uid_generator git:(main) ✗
```