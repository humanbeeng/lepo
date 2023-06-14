package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Client struct {
	Addr string
	Conn net.Conn
}

func New(loadbalanceraddr string) (*Client, error) {
	conn, err := net.Dial("tcp", loadbalanceraddr)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with %v", loadbalanceraddr)
	}
	return &Client{Conn: conn}, nil
}

func (c *Client) Register() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, byte(0))
	binary.Write(buf, binary.LittleEndian, int32(len(c.Addr)))
	binary.Write(buf, binary.LittleEndian, []byte(c.Addr))
	c.Conn.Write(buf.Bytes())
}

func (c *Client) DeRegister() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, byte(1))
	binary.Write(buf, binary.LittleEndian, int32(len(c.Addr)))
	binary.Write(buf, binary.LittleEndian, []byte(c.Addr))
	c.Conn.Write(buf.Bytes())
}
