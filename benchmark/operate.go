package benchmark

import (
	"cache_benchmark/cacheClient"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (b *BenchInfo)operate(id int, ch chan *result) {
	// total表示要压测的次数，默认值是1000
	// threads表示启动的线程数量，默认是1
	// count表示每个线程需要压测的次数
	count := b.Total/b.Threads
	client := cacheClient.New(b.Typ, b.Server)
	fmt.Println("打开连接,type is", b.Typ, "server is", b.Server)
	cmds := make([]*cacheClient.Cmd, 0)
	// vaelue以a为前缀,strings.Repeat表示将字符重复多少遍
	valuePrefix := strings.Repeat("a", b.ValueSize)
	r := &result{0,0,0, make([]statistic, 0)}
	for i:=0;i< count;i++ {
		//fmt.Println("id is", id, "total is", b.Total, "thread is", b.Threads, "keyspacelen is", b.Keyspacelen, "count is", count, "i is", i)
		// 取一个随机数据用来生成key
		var tmp int
		if b.Keyspacelen >0 {
			tmp = rand.Intn(b.Keyspacelen)
		} else {
			// 这个地方没理解，key的长度为0的话tmp等于啥，但是能保证key不重复
			tmp = id*count + i
		}
		// 以数字作为key
		key := fmt.Sprintf("%d", tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		name := b.Operation

		if b.Operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}
		//fmt.Println("id is", id, "i is", i, "key is", key, "name is", name)
		c := &cacheClient.Cmd{name, key, value, nil}
		fmt.Println("压测内容", c)
		if b.Pipelen > 1 {
			//fmt.Println("pipeline is running.................")
			// 如果是pipeline模式则将压测操作放到一个切片中
			cmds = append(cmds, c)
			// 如果切片的长度等于pipeline的长度了则执行压测
			if len(cmds) == b.Pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			fmt.Println("开始第", i, "次压测")
			run(client, c, r)
		}
	}
	if len(cmds) != 0 {
		pipeline(client, cmds, r)
	}
	ch <- r
}

func run(client cacheClient.Client, c *cacheClient.Cmd, r *result) {
	//fmt.Println("cmd is", c)
	// value
	expect := c.Value
	start := time.Now()
	client.Run(c)
	// 这是个求时间差值的好方法
	d := time.Now().Sub(start)
	//fmt.Println("d is", d)
	resultType := c.Name
	//fmt.Println("result type is", resultType)
	if resultType == "get" {
		//fmt.Println("value is", c.Value, "expect is", expect)
		if c.Value == "" {
			//fmt.Println("value is null")
			resultType = "miss"
		} else if c.Value != expect {
			//fmt.Println("value <> expect")
			panic(c)
		}
	}
	// 统计时间消耗
	r.addDuration(d, resultType)
}

func pipeline(client cacheClient.Client, cmds []*cacheClient.Cmd, r *result) {
	//fmt.Println("pipline function is running................")
	expect := make([]string, len(cmds))

	for i, c := range cmds {
		//fmt.Println("c is", c)
		if c.Name == "get" {
			// 这里为什么要把value放入到expect里呢？get操作的情况下value无意义
			expect[i] = c.Value
		}
	}
	//fmt.Println("expect is", expect)
	start := time.Now()
	// PipelinedRun该方法由于接收指针cmds作为参数，所以经过一番运作之后cmds里面就有内容了，我也经常用这种方法
	client.PipelinedRun(cmds)
	// 计算时间差
	d := time.Now().Sub(start)
	// 压测的具体操作已经执行完了，我们遍历其结果进行统计
	for i, c := range cmds {
		resultType := c.Name
		// get操作
		//fmt.Println("i is", i, "c is", c, "key is", c.Key, "value is", c.Value, "name is", c.Name)
		if resultType == "get" {
			// 判断value是否为空，如果唯恐则resultType = miss，这个不理解，怎么就成miss操作了呢
			if c.Value == "" {
				resultType = "miss"
			} else if c.Value != expect[i] {
				fmt.Println(expect[i])
				panic(c.Value)
			}
		}
		r.addDuration(d, resultType)
	}

}
