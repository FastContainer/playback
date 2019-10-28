package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
)

func main() {
	rcptAddr := "root@recipient"
	fromAddr := "root@sender"

	c, err := Dial("recipient:25")
	if err != nil {
		fmt.Printf("Dial\n%#v\n", err)
		return
	}
	err = c.Helo("localhost")
	if err != nil {
		return
	}
	err = c.MailFrom(fromAddr)
	if err != nil {
		return
	}
	err = c.RcptTo(rcptAddr)
	if err != nil {
		return
	}
	err = c.Data()
	if err != nil {
		return
	}

	body := "To: %s\r\n" +
		"Subject: %s\r\n" +
		"\r\n" +
		"This is the email body.\r\n"
	code, msg, err := c.Cmd(250, body, rcptAddr, "Hello, Gophers")
	if err != nil {
		fmt.Printf("HELO\n%d -- %s\n%#v", code, msg, err)
		return
	}
	return
}

type Client struct {
	conn *Conn
}

type Conn struct {
	textproto.Reader
	textproto.Writer
	textproto.Pipeline
	conn io.ReadWriteCloser
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) Cmd(format string, args ...interface{}) (id uint, err error) {
	id = c.Next()
	c.StartRequest(id)
	err = c.PrintfLine(format, args...)
	c.EndRequest(id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Dial(addr string) (*Client, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn := &Conn{
		Reader: textproto.Reader{R: bufio.NewReader(c)},
		Writer: textproto.Writer{W: bufio.NewWriter(c)},
		conn:   c,
	}
	_, _, err = conn.ReadResponse(220)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id := c.conn.Next()
	c.conn.StartRequest(id)
	err := c.conn.PrintfLine(format, args...)
	c.conn.EndRequest(id)
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
		fmt.Printf("HELO %s\nID: %d -- %s\n%#v\n", host, code, msg, err)
	}
	return err
}

func (c *Client) MailFrom(addr string) error {
	code, msg, err := c.Cmd(250, "MAIL FROM: "+addr)
	if err != nil {
		fmt.Printf("MAIL FROM: <%s>\nID: %d -- %s\n%#v\n", addr, code, msg, c.conn)
	}
	return err
}

func (c *Client) RcptTo(addr string) error {
	code, msg, err := c.Cmd(25, "RCPT TO: <%s>", addr)
	if err != nil {
		fmt.Printf("RCPT TO: <%s>\nID: %d -- %s\n%#v\n", addr, code, msg, c.conn)
	}
	return err
}

func (c *Client) Data() error {
	code, msg, err := c.Cmd(354, "DATA")
	if err != nil {
		fmt.Printf("DATA\nID: %d -- %s\n%#v\n", code, msg, c.conn)
	}
	return err
}
