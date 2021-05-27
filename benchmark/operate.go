package benchmark

import (
	"cache_benchmark/cacheClient"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (b *BenchInfo)operate(id int, ch chan *result) {
	count := b.Total/b.Threads
	client := cacheClient.New(b.Typ, b.Server)
	cmds := make([]*cacheClient.Cmd, 0)
	valuePrefix := strings.Repeat("a", b.ValueSize)
	r := &result{0,0,0, make([]statistic, 0)}
	for i:=0;i< count;i++ {
		var tmp int
		if b.Keyspacelen >0 {
			tmp = rand.Intn(b.Keyspacelen)
		} else {
			tmp = id*count + i
		}
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
		c := &cacheClient.Cmd{name, key, value, nil}
		if b.Pipelen > 1 {
			cmds = append(cmds, c)
			if len(cmds) == b.Pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			run(client, c, r)
		}
	}
	if len(cmds) != 0 {
		pipeline(client, cmds, r)
	}
	ch <- r
}

func run(client cacheClient.Client, c *cacheClient.Cmd, r *result) {
	expect := c.Value
	start := time.Now()
	client.Run(c)
	d := time.Now().Sub(start)
	resultType := c.Name
	if resultType == "get" {
		if c.Value == "" {
			resultType = "miss"
		} else if c.Value != expect {
			panic(c)
		}
	}
	r.addDuration(d, resultType)
}

func pipeline(client cacheClient.Client, cmds []*cacheClient.Cmd, r *result) {
	expect := make([]string, len(cmds))
	for i, c := range cmds {
		if c.Name == "get" {
			expect[i] = c.Value
		}
	}
	start := time.Now()
	client.PipelinedRun(cmds)
	d := time.Now().Sub(start)
	for i, c := range cmds {
		resultType := c.Name
		if resultType == "get" {
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
