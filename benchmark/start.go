package benchmark

import (
	"fmt"
	"time"
)

type BenchInfo struct {
	Typ, Server, Operation string
	Total, ValueSize, Threads, Keyspacelen, Pipelen int
}

func (b *BenchInfo) Start()  {
	// 声明一个result类型的channel，长度为threads
	ch := make(chan *result, b.Threads)
	// 初始化一个result类型的结构体实例
	// statistic结构体包含两个字段：count、time
	res := &result{0, 0,0, make([]statistic, 0)}
	// 开始计时
	start := time.Now()
	// 根据指定的线程数进行压测
	for i:=0;i<b.Threads;i++ {
		go b.operate(i, ch)
	}
	// 这个地方有点不理解了，为什么要并发读取结果呢，难道是担心读取结果的速度太慢影响压测效果？
	for i:=0;i<b.Threads;i++ {
		res.addResult(<- ch)
	}
	// 计时结束，这段时间是整个压测的过程，这个地方用sub来计算时间差的方式非常霸气
	d := time.Now().Sub(start)
	totalCount := res.getCount+res.missCount+res.setCount
	fmt.Printf("%d records get\n", res.getCount)
	fmt.Printf("%d records miss\n", res.missCount)
	fmt.Printf("%d records set\n", res.setCount)
	fmt.Printf("%f second total\n", d.Seconds())
	statCountSum := 0
	statTimeSum := time.Duration(0)
	for b, s := range res.statBuckets {
		if s.count == 0 {
			continue
		}
		statCountSum += s.count
		statTimeSum += s.time
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, b+1)
	}
	fmt.Printf("%d usec average for each reques\n", int64(statTimeSum/time.Second)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((res.getCount + res.setCount)*b.ValueSize)/1e6/d.Seconds())
	fmt.Println("rps is", float64(totalCount)/float64(d.Seconds()))
}
