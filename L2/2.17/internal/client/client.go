package client

import (
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	address string
	timeout time.Duration
	conn    net.Conn
}

func New(host, port string, timeout time.Duration) *Client {
	address := host + ":" + port
	return &Client{address: address, timeout: timeout, conn: nil}
}

func (client *Client) Run() error {
	conn, err := net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		io.Copy(conn, os.Stdin)
		conn.Close()
	}()

	io.Copy(os.Stdout, conn)

	return nil
}

func (c *Client) Connect() error { return nil }
func (c *Client) Close() error   { return nil }
