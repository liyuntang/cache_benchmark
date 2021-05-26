package cacheClient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *tcpClient)sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient)sendSet(key, value string)  {
	klen := len(key)
	vlen := len(value)
	c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (c *tcpClient)sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func readLen(r *bufio.Reader) int {
	tmp, err := r.ReadString(' ')
	if err != nil {
		log.Println(err)
		return 0
	}
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		log.Println(tmp, err)
		return 0
	}
	return l
}

func (c *tcpClient)recvResponse() (string, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		arr := make([]byte, -vlen)
		_, err := io.ReadFull(c.r, arr)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(arr))
	}
	value := make([]byte, vlen)
	_, err := io.ReadFull(c.r, value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (c *tcpClient)Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
		return
	}
	panic("unkonw cmd name "+ cmd.Name)
}

func (c *tcpClient)PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}

	for _, cmd := range cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value)
		}
		if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
}

func newTCPClient(server string) *tcpClient {
	c, err := net.Dial("tcp", server+":12346")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}







