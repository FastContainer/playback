package main

import (
	"fmt"
	"net/textproto"
)

func main() {
	run()
}

type Client struct {
	conn *textproto.Conn
}

// https://golang.org/src/net/smtp/smtp.go
// https://golang.org/src/net/textproto/textproto.go
func run() {
	client, err := Dial("recipient:25")
	if err != nil {
		fmt.Printf("Dial\n%#v\n", err)
		return
	}
	err := client.Helo("localhost")
	if err != nil {
		return
	}
	err := client.Mail("root@sender")
	if err != nil {
		return
	}

	text, err := client.conn.ReadDotBytes()
	if err != nil {
		fmt.Printf("ddd\n%#v\n", err)
		return
	}
	fmt.Printf("!!!!!!%#v\n", text)

	client.conn.ReadCodeLine(250)
	return
}

func Dial(addr string) (*Client, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	_, _, err = conn.ReadResponse(220)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.conn.Cmd(format, args)
	if err != nil {
		return 0, "", err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	code, msg, err := c.conn.ReadCodeLine(expectCode)
	return code, msg, err
}

func (c *Client) Helo(host string) error {
	code, msg, err := c.Cmd(250, "HELO %s", host)
	if err != nil {
		fmt.Printf("HELO\n%d -- %s\n%#v", code, msg, err)
		return
	}
	return err
}

func (c *Client) Mail(from string) error {
	code, msg, err := c.cmd(250, "MAIL FROM:<%s>", from)
	if err != nil {
		fmt.Printf("MAIL FROM\n%d -- %s\n%#v", code, msg, err)
		return
	}
	return err
}
