package main

import (
	"cache_benchmark/benchmark"
	"flag"
	"math/rand"
	"time"
)

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

func init()  {
	flag.StringVar(&typ, "type", "tcp", "cache server type,could be redis/http/tcp")
	flag.StringVar(&server, "h", "localhost", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "s", 1000, "data size of set/get value in bytes")
	flag.IntVar(&threads, "t", 1, "number of parallel connections")
	flag.StringVar(&operation, "o", "set", "test set, could be get/set/mixed")
	// key的长度，取得是随机数
	flag.IntVar(&keyspacelen, "k", 0, "key space len, use random keys from 0 to keyspacelen-1")
	flag.IntVar(&pipelen, "p", 1, "pipeline length")
	flag.Parse()
	/*
	如果不使用rand.Seed(seed int64)，每次运行，得到的随机数会一样，程序不停止，一直获取的随机数是不一样的；
	2、每次运行时rand.Seed(seed int64)，seed的值要不一样，这样生成的随机数才会和上次运行时生成的随机数不一样；
	3、rand.Intn(n int)得到的随机数int i，0 <= i < n。
	 */
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

//func main() {
//	server := flag.String("h", "localhost", "cache server address")
//	op := flag.String("c", "get", "command, could be get/set/del")
//	key := flag.String("k", "", "key")
//	value := flag.String("v", "", "value")
//	flag.Parse()
//	client := cacheClient.New("tcp", *server)
//	cmd := &cacheClient.Cmd{*op, *key, *value, nil}
//	client.Run(cmd)
//	if cmd.Error != nil {
//		fmt.Println("error:", cmd.Error)
//	} else {
//		fmt.Println(cmd.Value)
//	}
//}













