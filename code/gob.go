package code

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCode struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Code = (*GobCode)(nil)

func NewGobCode(conn io.ReadWriteCloser) Code {
	buf := bufio.NewWriter(conn)
	return &GobCode{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCode) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCode) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCode) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCode) Close() error {
	return c.conn.Close()
}
