package healthcheck

import (
	"net"
	"time"
)

// TCPChecker defines a TCP checker.
type TCPChecker struct {
	Timeout time.Duration
}

// CheckSlice checks a slice of addresses.
func (t *TCPChecker) CheckSlice(addresses []string, result ResultChan) {
	for _, address := range addresses {
		go func(ad string) {
			if err := t.Check(ad); err == nil {
				result.OKChan <- ResultItem{Address: ad}
			} else {
				result.ErrorChan <- ResultItem{Address: ad, Error: err}
			}
		}(address)
	}
}

// Check checks a single address
func (t *TCPChecker) Check(address string) error {
	conn, err := net.DialTimeout("tcp", address, t.Timeout)
	if err == nil {
		_ = conn.Close()
	}

	return err
}
