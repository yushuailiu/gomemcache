package gomemcache

import (
	"sync"
	"net"
	"strings"
	"hash/crc32"
)

// ServerList is a simple ServerSelector. Its zero value is usable.
type ServerList struct {
	sync.RWMutex
	addrs []net.Addr
}

type staticAddr struct {
	ntw, str string
}

func (s *staticAddr) Network() string { return s.ntw }

func (s *staticAddr) String() string  { return s.str }

func newStaticAddr(a net.Addr) net.Addr {
	return &staticAddr{
		ntw: a.Network(),
		str: a.String(),
	}
}

// SetServers changes a ServerList's set of servers at runtime and is safe
// for concurrent use by multiple goroutines.
//
// Each server is given equal weight. A server is given more weight
// if it's listed multiple times.
//
// SetServers returns an error if any of the server names fail to resolve.
// No attempt is made to connect to the server. If any error is returned,
// no changes are made to the ServerList.
func (ss *ServerList) SetServers(servers ...string) error {
	naddr := make([]net.Addr, len(servers))

	for i, server := range servers {
		if strings.Contains(server, "/") {
			addr, err := net.ResolveUnixAddr("unix", server)
			if err != nil {
				return err
			}
			naddr[i] = newStaticAddr(addr)
		} else {
			tcpaddr, err := net.ResolveTCPAddr("tcp", server)
			if err != nil {
				return err
			}
			naddr[i] = newStaticAddr(tcpaddr)
		}
	}

	ss.Lock()
	defer ss.Unlock()
	ss.addrs = naddr
	return nil
}

// Each iterates over each server calling the given function
func (ss *ServerList) Each(f func(net.Addr) error) error {
	ss.Lock()
	defer ss.Unlock()

	for _, a := range ss.addrs {
		if err := f(a); nil != err {
			return err
		}
	}
	return nil
}

var keyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 256)
		return &b;
	},
}

func (ss *ServerList) PickServer(key string) (net.Addr, error) {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.addrs) == 0 {
		return nil, ErrNoServers
	}

	bufp := keyBufPool.Get().(*[]byte)
	n := copy(*bufp, key)
	cs := crc32.ChecksumIEEE((*bufp)[:n])
	keyBufPool.Put(bufp)

	return ss.addrs[cs%uint32(len(ss.addrs))], nil
}