package gomemcache

import (
	"net"
	"errors"
)

var (
	// ErrCacheMiss means that a Get failed because the item wasn't present.
	ErrCacheMiss = errors.New("memcache: cache miss")

	// ErrNoServers is returned when no servers are configured or available.
	ErrNoServers = errors.New("memcache: no servers configured or available")

	ErrCASConflict = errors.New("memcache: compare-and-swap conflict")

	ErrNotStored = errors.New("memecache: item not stored")

	ErrMalformedKey = errors.New("malformed: key is too long or contains invalid characters")
)

type ConnectTimeoutError struct {
	Addr net.Addr
}

func (cte *ConnectTimeoutError) Error() string {
	return "memcache: connect timeout to " + cte.Addr.String()
}