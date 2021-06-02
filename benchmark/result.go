package benchmark

import "time"

type statistic struct {
	count int
	time time.Duration
}

type result struct {
	getCount	int
	missCount	int
	setCount	int
	statBuckets	[]statistic
}

func (r *result)addStatistic(bucket int, stat statistic) {
	if bucket > len(r.statBuckets)-1 {
		newStatBuckets := make([]statistic, bucket+1)
		copy(newStatBuckets, r.statBuckets)
		r.statBuckets = newStatBuckets
	}
	s :=  r.statBuckets[bucket]
	s.count += stat.count
	s.time += stat.time
	r.statBuckets[bucket] = s
}

func (r *result)addDuration(d time.Duration, typ string) {
	bucket := int(d /time.Millisecond)
	r.addStatistic(bucket, statistic{1, d})
	if typ == "get" {
		r.getCount++
	}else if typ == "set" {
		r.setCount++
	} else {
		r.missCount++
	}
}

// 统计压测信息的方法
func (r *result)addResult(src *result) {
	// statBuckets是result结构体的一个字段，是statist类型的一个切片
	// 这里遍历statist切片，并对切片内的数据进行处理
	for b, s := range src.statBuckets {
		// b是下标号，一个自然数，s是statistic结构体
		r.addStatistic(b, s)
	}
	r.getCount += src.getCount
	r.missCount += src.missCount
	r.setCount += src.setCount
}













