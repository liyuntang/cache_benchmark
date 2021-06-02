package cacheClient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient)get(key string) string {
	resp, err := c.Get(c.server+key)
	if err != nil {
		log.Println(key)
		panic(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (c *httpClient)set(key, value string) {
	url := fmt.Sprintf("%s/%s", c.server, key)
	fmt.Println("key is", key, "value is", string(value), "url is", url)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(value))
	if err != nil {
		log.Println("new reques is bad",key)
		panic(err)
	}
	defer req.Body.Close()
	resp, err := c.Do(req)
	if err != nil {
		log.Println("do is bad",key)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func (c *httpClient)Run(cmd *Cmd)  {
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		fmt.Println("http set is running...................")
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unkonw cmd name "+ cmd.Name)
}

func newHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://"+server+":12345/cache"}
}

func (c *httpClient)Close()  {
	c.Close()
}
func (c *httpClient)PipelinedRun([]*Cmd) {
	panic("httpClient pipelined run ont implement")
}
