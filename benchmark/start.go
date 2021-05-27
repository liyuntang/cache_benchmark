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
	ch := make(chan *result, b.Threads)
	res := &result{0, 0,0, make([]statistic, 0)}
	start := time.Now()
	for i:=0;i<b.Threads;i++ {
		go b.operate(i, ch)
	}
	for i:=0;i<b.Threads;i++ {
		res.addResult(<- ch)
	}
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
