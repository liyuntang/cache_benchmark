package cacheClient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	tmpReq := fmt.Sprintf("G%d %s", klen, key)
	fmt.Println(tmpReq, c)
	c.Write([]byte(tmpReq))
}

func (c *tcpClient)sendSet(key, value string)  {
	klen := len(key)
	vlen := len(value)
	requestData := fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)
	c.Write([]byte(requestData))
}

func (c *tcpClient)sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func readLen(r *bufio.Reader) int {
	fmt.Println("read len go.........................")
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("read all is bad, err is", err)
	} else {
		fmt.Println("buf is", string(buf))
	}
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
	fmt.Println("tmp is", tmp, "l is", l)
	return l
}

// 接收对面tcp服务返回的信息
func (c *tcpClient)recvResponse() (string, error) {
	fmt.Println("go up>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	vlen := readLen(c.r)
	fmt.Println("vlen is", vlen)
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
		//fmt.Println("tcp get is running...............", cmd, cmd.Key)
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		return
	}
	if cmd.Name == "set" {
		fmt.Println("tcpclient set operation is running...................")
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
	endPoint := fmt.Sprintf("%s:12346", server)
	//address := fmt.Printf("%s:12346", server)
	//fmt.Println("new tcp client...............", endPoint)
	c, err := net.Dial("tcp", endPoint)
	if err != nil {
		panic(err)
	}
	//fmt.Println("new tcp client is ok", endPoint)
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}







