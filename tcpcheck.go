package healthcheck

import (
	"net"
	"time"
)

type TcpChecker struct {
	Timeout time.Duration
}

func (t *TcpChecker) TcpCheckSlice(addresses []string, result ResultChan) {
	for _, address := range addresses {
		go func(ad string) {
			if err := t.TcpCheck(ad); err == nil {
				result.OKChan <- ResultItem{Address: ad}
			} else {
				result.ErrorChan <- ResultItem{Address: ad, Error: err}
			}
		}(address)
	}
}

func (t *TcpChecker) TcpCheck(address string) error {
	conn, err := net.DialTimeout("tcp", address, t.Timeout)
	if err == nil {
		_ = conn.Close()
	}

	return err
}
