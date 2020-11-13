package main

import (
	"flag"
	"fmt"
	"github.com/bingoohuang/healthcheck"
	"strings"
	"time"

	"github.com/bingoohuang/gometrics/metric"
)

func main() {
	ap := app{addrs: healthcheck.NewAddresses()}
	ap.parseFlags()
	ap.goTCPCheck()
}

type app struct {
	dialTimeout time.Duration

	addrs    healthcheck.Addresses
	result   healthcheck.ResultChan
	interval time.Duration
}

func (a *app) parseFlags() {
	addrListFileName := flag.String("f", "", "address file name(ip:port per line)")
	addrList := flag.String("a", "", "address list (format ip:port ip:port)")
	timeout := flag.String("t", "3s", "connect timeout(default 3s)")
	interval := flag.String("interval", "", "checking interval")
	flag.Parse()

	a.dialTimeout = healthcheck.MustParseDuration(*timeout)
	a.interval = healthcheck.MustParseDuration(*interval)
	a.addrs.ParseAddress(*addrListFileName, *addrList)
}

func (a *app) goTCPCheck() {
	checker := healthcheck.TCPChecker{Timeout: a.dialTimeout}

	if a.interval > 0 {
		a.intervalCheck(checker)
	} else {
		a.onceCheck(checker)
	}
}

// intervalCheck 定期检查，1.生成metrics 2. 打印输出结果
func (a *app) intervalCheck(checker healthcheck.TCPChecker) {
	if len(a.addrs) == 0 {
		fmt.Println("no addresses specified")
		return
	}

	metric.Start()
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		for _, addr := range a.addrs {
			keys := []string{DefaultTo(addr.Desc, "HealthCheck")}
			addr_ := strings.ReplaceAll(addr.Addr, ".", "_")
			keys = append(keys, strings.SplitN(addr_, ":", 2)...)

			qps := metric.QPS(keys...)
			succRate := metric.SuccessRate(keys...)
			failRate := metric.FailRate(keys...)

			err := checker.Check(addr.Addr)
			if err == nil {
				succRate.IncrSuccess()
			} else {
				failRate.IncrFail()
			}

			qps.Record(1)
			succRate.IncrTotal()
			failRate.IncrTotal()
		}

		<-ticker.C
	}
}

func DefaultTo(a, b string) string {
	if a == "" {
		return b
	}

	return a
}

// onceCheck 一次性检查，并且打印输出结果
func (a *app) onceCheck(checker healthcheck.TCPChecker) {
	a.result = healthcheck.NewResult()
	checker.CheckSlice(a.addrs, a.result)

	results := a.result.WaitResults(a.addrs.Len())
	results.PrintResult()
}
