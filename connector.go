// Package dialer provides an abstraction on the
// OS network dialer to conveniently use in its
// place an object accessing data in memory
package dialer

import (
	"context"
	"net"
)

// Connector abstracts all connectors
type Connector interface {
	Connect() (net.Conn, error)
}

// OSConn implements Connector using the net package
type OSConn struct {
	Dlr  *net.Dialer
	Net  string
	Addr string
	// This field may be nil
	Ctx context.Context
}

// Connect is the Connector implementation
func (c *OSConn) Connect() (r net.Conn, e error) {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
	r, e = c.Dlr.DialContext(c.Ctx, c.Net, c.Addr)
	return
}
