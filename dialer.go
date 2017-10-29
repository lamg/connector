// Package dialer provides an abstraction on the
// OS network dialer to conveniently use in its
// place an object accessing data in memory
package dialer

import (
	"net"
)

// Dialer is an interface for treating a dialer
// returning a network connection and a dialer returning
// a connection with its data in memory as the same
type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

// OSDialer implementation of Dialer using OS network
type OSDialer struct {
	Dlr *net.Dialer
}

// NewOSDialer creates a new OSDialer
// local: local address for making the connection
func NewOSDialer(local net.Addr) (s *OSDialer) {
	dlr := &net.Dialer{LocalAddr: local}
	s = &OSDialer{Dlr: dlr}
	return
}

// Dial dials using the supplied net.IP as local
// address
func (s *OSDialer) Dial(nt, addr string) (c net.Conn, e error) {
	c, e = s.Dlr.Dial(nt, addr)
	return
}
