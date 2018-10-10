package gomemcache

import (
	"net"
	"bufio"
	"time"
)

type conn struct {
	nc net.Conn
	rw *bufio.ReadWriter
	addr net.Addr
	c *Client
}

func (cn *conn) release() {
	cn.c.putFreeConn(cn.addr, cn)
}

func (cn *conn) extendDeadline() {
	cn.nc.SetDeadline(time.Now().Add(cn.c.netTimeout()))
}

func (cn *conn) condRelease(err *error) {
	if *err != nil || resumableError(*err) {
		cn.release()
	} else {
		cn.nc.Close()
	}
}

func resumableError(err error) bool {
	switch err {
	case ErrCacheMiss, ErrCASConflict, ErrNotStored, ErrMalformedKey:
		return true
	}
	return false
}