# Snowflake
[![GoDoc](https://godoc.org/github.com/SaurabhGoyal/Snowflake?status.svg)](https://godoc.org/github.com/SaurabhGoyal/Snowflake) [![Go Report Card](https://goreportcard.com/badge/github.com/SaurabhGoyal/Snowflake)](https://goreportcard.com/report/github.com/SaurabhGoyal/Snowflake)

Snowflake is a go package that provides a simple implementation of unique-id generation based on [Twitter's snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) logic.
```
+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
+--------------------------------------------------------------------------+
```

## Status
Project has been tested and benchmarked but not battle tested in production environment yet.

## Logic
- There are two components - `Generator` and `Config`.
- `Config` has a default implementation which provides 2 ^ 21 (2097152) unique ids per millisecond (1024 per node, across 2048 nodes) till year 2149 (~139 years from epoch). These values can be customised to support higher throughput or longer period of valid generation of unique IDs. Default config uses following details -
  - epoch - `2010/12/12/23/59/59/0 UTC`.
  - timestampBits - 42
  - nodeIdBits - 11

## Hands On
- Install (Assuming you already have a working Go environment, if not please see [this page](https://go.dev/doc/install) first)
```
go get github.com/SaurabhGoyal/Snowflake
```
- Use
https://go.dev/play/p/GeX6m9UKkiD
```go
package main

import (
	"log"
	"time"

	snowflake "github.com/SaurabhGoyal/Snowflake/snowflake"
)

func main() {

	// With default config
	config, _ := snowflake.InitDefaultConfig()
	uidGen, _ := snowflake.InitGenerator(config, 1)
	uid, err := uidGen.Get()
	if err != nil {
		panic(err)
	}
	log.Printf("UID from default config - %d", uid)

	// With custom config (Ex.- less lifetime, more servers, less throughput per server)
	epoch := uint64(time.Date(2015, 12, 12, 23, 59, 59, 0, time.UTC).UnixMilli())
	config, _ = snowflake.InitConfig(epoch, 41, 14)
	uidGen, _ = snowflake.InitGenerator(config, 1)
	uid, err = uidGen.Get()
	if err != nil {
		panic(err)
	}
	log.Printf("UID from default config - %d", uid)

}
```
Output
```
➜  Snowflake git:(main) ✗ go run main.go
2022/09/20 11:31:51
==============================================================================
Initialising Snowflake Unique ID Generator Config
==============================================================================
Config (64 bit ID)
+------------------------------------------------------------------------+
| 1 Bit Unused | 42 Bit Timestamp |  11 Bit NodeID  | 10 Bit Sequence ID |
+------------------------------------------------------------------------+
Output (TP = requests per millisecond)
+---------------------------------------------------------------------------------+
| 139.46 Years of uniqueness lifetime | 2097152 TP across 2048 servers  | 1024 TP per server |
+---------------------------------------------------------------------------------+
==============================================================================

2022/09/20 11:31:51 UID from default config - 778998251375297537
2022/09/20 11:31:51
==============================================================================
Initialising Snowflake Unique ID Generator Config
==============================================================================
Config (64 bit ID)
+------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  14 Bit NodeID  | 8 Bit Sequence ID |
+------------------------------------------------------------------------+
Output (TP = requests per millisecond)
+---------------------------------------------------------------------------------+
| 69.73 Years of uniqueness lifetime | 4194304 TP across 16384 servers  | 256 TP per server |
+---------------------------------------------------------------------------------+
==============================================================================

2022/09/20 11:31:51 UID from default config - 896276260164993281
➜  Snowflake git:(main) ✗
```

## Benchmarks
- Machine-1
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

- Machine-2
```
➜  Snowflake git:(main) ✗ go test -v -run=^# -bench=. -count=4 ./... -benchmem
?       github.com/SaurabhGoyal/Snowflake       [no test files]
goos: linux
goarch: amd64
pkg: github.com/SaurabhGoyal/Snowflake/snowflake
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkGet
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8            1000000              1173 ns/op               0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8             953890              1184 ns/op               0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8             949268              1195 ns/op               0 B/op          0 allocs/op
BenchmarkGet/Default_config_-_Tuned_for_high_distribution_of_nodes_(11_bits)_and_moderate_throughput_(10_bits)_per_node-8             991071              1159 ns/op               0 B/op          0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                  14161412                95.48 ns/op            0 B/op          0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                  13399654                95.54 ns/op            0 B/op          0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                  13667899                94.33 ns/op            0 B/op          0 allocs/op
BenchmarkGet/Custom_config_-_Tuned_for_low_distribution_of_nodes_(7_bits)_and_high_throughput_(14_bits)_per_node-8                  13927473                95.93 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/SaurabhGoyal/Snowflake/snowflake     10.261s
?       github.com/SaurabhGoyal/Snowflake/uid   [no test files]
➜  Snowflake git:(main) ✗
```