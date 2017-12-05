// Package nettest provides an abstraction on the
// OS network dialer to conveniently use in its
// place an object accessing data in memory
package nettest

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"
)

// Connector abstracts all connectors
type Connector interface {
	// SetDialer and SetContext are optional calls before
	// Connect call, all parameters may be nil
	SetDialer(*net.Dialer)
	SetContext(context.Context)
	Connect(string, string) (net.Conn, error)
}

// OSConn implements Connector using the net package
type OSConn struct {
	dlr *net.Dialer
	ctx context.Context
}

// SetContext is part of Connector implementation
func (c *OSConn) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// SetDialer is part of Connector implementation
func (c *OSConn) SetDialer(d *net.Dialer) {
	c.dlr = d
}

// Connect is part of Connector implementation
func (c *OSConn) Connect(n, a string) (r net.Conn, e error) {
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.dlr == nil {
		c.dlr = new(net.Dialer)
	}
	r, e = c.dlr.DialContext(c.ctx, n, a)
	return
}

// MemConn is a Connector implementation using strings
// stored in memory as content for connections
type MemConn struct {
	Data  map[string]string
	LAddr net.Addr
	dlr   *net.Dialer
	ctx   context.Context
}

// SetDialer is part of Connector implementation
func (c *MemConn) SetDialer(d *net.Dialer) {
	c.dlr = d
}

// SetContext is part of Connector implementation
func (c *MemConn) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// Connect is part of Connector implementation
func (c *MemConn) Connect(n, a string) (r net.Conn, e error) {
	s, ok := c.Data[a]
	if !ok {
		e = fmt.Errorf("Host %s not found", a)
	}
	if c.dlr == nil {
		c.dlr = &net.Dialer{
			LocalAddr: c.LAddr,
		}
	}
	if e == nil {
		r = &MConn{
			Buffer:   bytes.NewBufferString(s),
			LAddr:    c.dlr.LocalAddr,
			RAddr:    &net.IPAddr{IP: net.ParseIP(a)},
			DeadLine: c.dlr.Deadline,
		}
	}
	return
}

// MConn implements net.Conn
type MConn struct {
	Buffer   *bytes.Buffer
	RAddr    net.Addr
	LAddr    net.Addr
	ReadDl   time.Time
	WriteDl  time.Time
	DeadLine time.Time
	closed   bool
}

func (c *MConn) Read(p []byte) (n int, e error) {
	if !c.closed {
		n, e = c.Buffer.Read(p)
	} else {
		n, e = 0, fmt.Errorf("Cannot read on closed connection")
	}
	return
}

func (c *MConn) Write(p []byte) (n int, e error) {
	if !c.closed {
		n, e = c.Buffer.Write(p)
	} else {
		n, e = 0, fmt.Errorf("Cannot write on closed connection")
	}
	return
}

// SetDeadline is part of net.Conn implementation
func (c *MConn) SetDeadline(t time.Time) (e error) {
	c.DeadLine = t
	return
}

// SetWriteDeadline is part of net.Conn implementation
func (c *MConn) SetWriteDeadline(t time.Time) (e error) {
	c.WriteDl = t
	return
}

// SetReadDeadline is part of net.Conn implementation
func (c *MConn) SetReadDeadline(t time.Time) (e error) {
	c.ReadDl = t
	return
}

// LocalAddr is part of net.Conn implementation
func (c *MConn) LocalAddr() (a net.Addr) {
	a = c.LAddr
	return
}

// RemoteAddr is part of net.Conn implementation
func (c *MConn) RemoteAddr() (a net.Addr) {
	a = c.RAddr
	return
}

// Close is part of net.Conn implementation
func (c *MConn) Close() (e error) {
	if !c.closed {
		c.closed = true
	}
	return
}
