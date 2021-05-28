package cacheClient

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type redisClient struct {
	*redis.Client
}

func (r *redisClient)get(key string) (string, error) {
	res, err := r.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func (r *redisClient) set(key, value string) error {
	return r.Set(ctx, key, value, 0).Err()
}

func (r *redisClient)del(key string)error {
	return r.Del(ctx, key).Err()
}

func (r *redisClient)Run(c *Cmd) {
	if c.Name == "get" {
		c.Value, c.Error = r.get(c.Key)
		return
	}
	if c.Name == "set" {
		c.Error = r.set(c.Key, c.Value)
		return
	}
	if c.Name == "del" {
		c.Error = r.del(c.Key)
		return
	}
	panic("unkonw cmd name "+ c.Name)
}

func (r *redisClient)PipelinedRun(cmds []*Cmd) {
	// 如果cmds的长度为0则表示不进行压测，直退出程序
	if len(cmds) == 0 {
		return
	}
	// 开启一个pipeline模式，之所以能够开启pipeline是因为r是redis.clent结构体
	pipe := r.Pipeline()
	// 这个地方对pipeline的设计不太理解，记住这个使用姿势就行了
	cmders := make([]redis.Cmder, len(cmds))
	for i, c := range cmds {
		//fmt.Println("c is", c)
		// 根据operation的不同将压测内容放入不同的pipeline队列中，目前支持get、set、delete三个操作
		if c.Name == "get" {
			cmders[i] = pipe.Get(ctx, c.Key)
		} else if c.Name == "set" {
			cmders[i] = pipe.Set(ctx, c.Key, c.Value, 0)
		} else if c.Name == "del" {
			cmders[i] = pipe.Del(ctx, c.Key)
		} else {
			// 这个地方感觉可以不用写这个了，get、set、del之外的操作直接丢弃即可
			panic("unknow cmd name "+c.Name)
		}
	}
	// 执行pipeline操作
	_, err := pipe.Exec(ctx)
	// 这里的redis.nil是个什么鬼
	if err != nil && err != redis.Nil {
		panic(err)
	}
	// 这里又来了一遍，啥目的
	for i, c := range cmds {
		//fmt.Println("i is", i, "c is", c)
		if c.Name == "get" {

			// 这是个什么逻辑，exec之后不是执行完了吗？不是已经拿到结果了吗，为什么还要redis.stringcmd，卧槽了，这个地方可能是排序，将request与response一一对应
			value, e := cmders[i].(*redis.StringCmd).Result()
			//fmt.Println("走起", "value is", value, "e is", e)
			// 卧槽，懂了，redis.nil表示get或del操作时如果redis中没有该key则返回nil
			if e != redis.Nil {
				// 表示该key存在，key存在的话将value设置为kong
				value, e = "", nil
			}
			// 如果err等于redis.nil则
			c.Value, c.Error = value, e
		} else {
			// 这说明执行的是set或del操作，为什么直接取err了呢
			c.Error = cmders[i].Err()
		}
	}
}

func newRedisClient(server string) *redisClient {
	return &redisClient{redis.NewClient(&redis.Options{Addr: server+":6379", ReadTimeout: -1})}
}






