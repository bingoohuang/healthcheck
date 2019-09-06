package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

type Result struct {
	Addr string
	Err  error
}

type app struct {
	addrListFileName string
	addrList         string
	dialTimeout      time.Duration

	addrs []string

	okChan, errChan chan Result
	oks, errs       []Result
}

func main() {
	ap := app{}
	ap.parseFlags()
	ap.prepareAddrs()
	ap.goTcpCheck()
	ap.waitResults()
	ap.printResult()
}

func (a *app) parseFlags() {
	addrListFileName := flag.String("f", "", "address file name(ip:port per line)")
	addrList := flag.String("a", "", "address list (format ip:port ip:port)")
	timeout := flag.String("t", "3s", "connect timeout(default 3s)")
	flag.Parse()

	dialTimeout, err := time.ParseDuration(*timeout)
	if err != nil {
		panic(err)
	}

	a.addrListFileName = *addrListFileName
	a.addrList = *addrList
	a.dialTimeout = dialTimeout
}

func (a *app) prepareAddrs() {
	a.addrs = make([]string, 0)

	if a.addrListFileName != "" {
		expanded, _ := homedir.Expand(a.addrListFileName)
		b, err := ioutil.ReadFile(expanded)
		if err != nil {
			panic(err)
		}

		lines := strings.Split(string(b), "\n")
		a.addAddr(lines)
	}

	if a.addrList != "" {
		lines := strings.Fields(a.addrList)
		a.addAddr(lines)
	}
}

func (a *app) goTcpCheck() {
	a.okChan = make(chan Result)
	a.errChan = make(chan Result)
	for _, addr := range a.addrs {
		go a.tcpCheck(addr)
	}
}

func (a *app) printResult() {
	total := len(a.oks) + len(a.errs)
	if len(a.errs) > 0 {
		fmt.Printf("Failed %d/%d:\n", len(a.errs), total)
		for _, addr := range a.errs {
			fmt.Printf("%s, error %v\n", addr.Addr, addr.Err)
		}
	}

	if len(a.oks) > 0 {
		if len(a.errs) > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("OK %d/%d:\n", len(a.oks), total)
		for _, ok := range a.oks {
			fmt.Printf("%s\n", ok.Addr)
		}
	}
}

func (a *app) waitResults() {
	a.oks = make([]Result, 0)
	a.errs = make([]Result, 0)
	cnt := len(a.addrs)
	for cnt > 0 {
		select {
		case ok := <-a.okChan:
			a.oks = append(a.oks, ok)
			cnt--
		case err := <-a.errChan:
			a.errs = append(a.errs, err)
			cnt--
		}
	}
}

func (a *app) addAddr(lines []string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if contains(a.addrs, line) {
			fmt.Printf("duplicate address %s\n", line)
		} else {
			a.addrs = append(a.addrs, line)
		}
	}
}

func (a *app) tcpCheck(addr string) {
	conn, err := net.DialTimeout("tcp", addr, a.dialTimeout)
	if err == nil {
		_ = conn.Close()
		a.okChan <- Result{addr, err}
	} else {
		a.errChan <- Result{addr, err}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if item == s {
			return true
		}
	}

	return false
}
