package healthcheck_test

import (
	"github.com/bingoohuang/healthcheck"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAddress(t *testing.T) {
	var addrs healthcheck.Addresses
	addrs.ParseAddress("testdata/addr.txt", "192.168.126.182:22")
	assert.Equal(t, healthcheck.Addresses([]healthcheck.Address{
		{
			Addr: "192.168.126.16:22",
			Desc: "测试机1",
		}, {
			Addr: "192.168.126.18:22",
			Desc: "测试机2",
		}, {
			Addr: "192.168.126.182:22",
			Desc: "",
		},
	}), addrs)
}
