package main

import (
	"cache_benchmark/benchmark"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

func init()  {
	flag.StringVar(&typ, "type", "redis", "cache server type")
	flag.StringVar(&server, "h", "localhost", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "s", 1000, "data size of set/get value in bytes")
	flag.IntVar(&threads, "t", 1, "number of parallel connections")
	flag.StringVar(&operation, "o", "set", "test set, could be get/set/mixed")
	flag.IntVar(&keyspacelen, "k", 0, "key space len, use random keys from 0 to keyspacelen-1")
	flag.IntVar(&pipelen, "p", 1, "pipeline length")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("server is", server)
	fmt.Println("total", total, "requests")
	fmt.Println("data size is", valueSize)
	fmt.Println("we have", threads, "connections")
	fmt.Println("operation is", operation)
	fmt.Println("keyspacelen is", keyspacelen)
	fmt.Println("pipeline length is", pipelen)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	info := &benchmark.BenchInfo{}
	info.Typ = typ
	info.Server = server
	info.Total = total
	info.ValueSize = valueSize
	info.Threads = threads
	info.Operation = operation
	info.Keyspacelen = keyspacelen
	info.Pipelen = pipelen
	info.Start()
}














