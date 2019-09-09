package main

import (
	"flag"
	"time"

	"github.com/bingoohuang/healthcheck"
)

func main() {
	ap := app{addrs: healthcheck.NewAddresses()}
	ap.parseFlags()
	ap.goTcpCheck()
}

type app struct {
	dialTimeout time.Duration

	addrs  healthcheck.Addresses
	result healthcheck.ResultChan
}

func (a *app) parseFlags() {
	addrListFileName := flag.String("f", "", "address file name(ip:port per line)")
	addrList := flag.String("a", "", "address list (format ip:port ip:port)")
	timeout := flag.String("t", "3s", "connect timeout(default 3s)")
	flag.Parse()

	a.dialTimeout = healthcheck.MustParseDuration(*timeout)
	a.addrs.PrepareAddress(*addrListFileName, *addrList)
}

func (a *app) goTcpCheck() {
	checker := healthcheck.TcpChecker{Timeout: a.dialTimeout}

	a.result = healthcheck.NewResult()
	checker.TcpCheckSlice(a.addrs, a.result)

	results := a.result.WaitResults(a.addrs.Len())
	results.PrintResult()
}
