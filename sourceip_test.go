package search

import (
	"context"
	"net"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	t.Run("Given valid IPv4 and IPv6", func(t *testing.T) {
		var tests = []struct {
			sourceIP SourceIP
			expected net.IP
		}{
			{sourceIP: NewSourceIP("172.0.0.1"), expected: net.ParseIP("172.0.0.1")},
			{sourceIP: NewSourceIP("10.0.0.0"), expected: net.ParseIP("10.0.0.0")},
			{sourceIP: NewSourceIP("192.168.0.1"), expected: net.ParseIP("192.168.0.1")},
			{sourceIP: NewSourceIP("2001:db8::68"), expected: net.ParseIP("2001:db8::68")},
			{sourceIP: NewSourceIP("2001:db8:85a3:0:0:8a2e:370:7334"), expected: net.ParseIP("2001:db8:85a3:0:0:8a2e:370:7334")},
			{sourceIP: NewSourceIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), expected: net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")},
		}

		for _, test := range tests {
			ctx, err := test.sourceIP.NewContext(context.Background())
			if err != nil {
				t.Fatal("Unexpected error: ", err)
			}

			var actual SourceIP
			if err := actual.FromContext(ctx); err != nil {
				t.Fatal("Unexpected error: ", err)
			}

			if !reflect.DeepEqual(test.expected, net.IP(actual)) {
				t.Errorf("Mismatch IP address.\nBad test case: %+v\nExpected %v, but got %v", test, test.expected, actual)
			}
		}
	})

	t.Run("Given invalid IP", func(t *testing.T) {
		var tests = []struct {
			sourceIP SourceIP
		}{
			{sourceIP: NewSourceIP("")},
			{sourceIP: NewSourceIP("test.example")},
			{sourceIP: NewSourceIP("172.0.0.1:7000")},
		}

		for _, test := range tests {
			s := test.sourceIP
			_, err := s.NewContext(context.Background())
			if err == nil {
				t.Fatal("Expected error didn't occur. Shouldn't be able to create context.")
			}
		}
	})
}
