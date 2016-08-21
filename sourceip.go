package search

import (
	"context"
	"fmt"
	"net"
)

// KeyIPAddr is the key used to identify user IP address in a context.
const KeyIPAddr = iota

// SourceIP is the IP address that originated the search query.
type SourceIP net.IP

// NewSourceIP returns a new instance of SourceIP.
func NewSourceIP(ipAddr string) SourceIP {
	return SourceIP(net.ParseIP(ipAddr))
}

// NewContext creates a copy of ctx containing s as the source IP address.
// It returns an error if s is empty.
func (s *SourceIP) NewContext(ctx context.Context) (context.Context, error) {
	if *s == nil {
		return nil, fmt.Errorf("Can't create context with empty source IP")
	}
	return context.WithValue(ctx, KeyIPAddr, *s), nil
}

// FromContext retrieves the source IP address from ctx.
// It returns an error if type assertion failed.
func (s *SourceIP) FromContext(ctx context.Context) error {
	sourceIP, ok := ctx.Value(KeyIPAddr).(SourceIP)
	if !ok {
		return fmt.Errorf("Unable to extract IP address from context")
	}

	*s = sourceIP
	return nil
}
