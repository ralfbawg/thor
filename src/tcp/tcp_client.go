package tcp

import (
	"github.com/panjf2000/ants"
	"util"
	"bytes"
	"common/logging"
	"io"
)

var (
	tcpCPoolExtendFactor       = 0.8
	tcpCPoolDefaultSize        = 10000
	tcpCPool, _                = ants.NewPool(tcpCPoolDefaultSize)
	funcs                      = make([]func(), 0)
	newline                    = []byte{'\n'}
	space                      = []byte{' '}
	bindClients, unbindClients = util.NewConcMap(), util.NewConcMap()
)

func (c *TcpClient) run() {
	util.SubmitTaskAndResize(tcpCPool, tcpCPoolDefaultSize, tcpCPoolExtendFactor, append(funcs[:0], c.Write, c.Read))
}

func (c *TcpClient) Write() {
	for {
		select {
		case msg := <-c.send:
			c.conn.Write(msg)
		}
	}
}
func (c *TcpClient) Read() {

	for {
		b := make([]byte, 256)
		n, err := c.conn.Read(b)
		if n != 0 {
			b = bytes.TrimSpace(b)
			bs := make([][]byte, 1)
			if bytes.Contains(b, newline) {
				bs = bytes.Split(b, newline)
			} else {
				bs = append(bs, b)
			}
			for _, v := range bs {
				b = v[0:n]
				logging.Info("get message %s from %s", string(b), c.ip)
				ProcessTcpMsg(b)
			}
			b = b[:0]
		} else if err == io.EOF {
			logging.Info("got %v; want %v", err, io.EOF)
			continue
		} else {
			logging.Info("got %v; want %v", err, io.EOF)
			break
		}

	}

}
