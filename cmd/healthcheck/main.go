package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type Result struct {
	Addr string
	Err  error
}

func main() {
	addrListFileName := flag.String("f", "", "address file name(ip:port per line)")
	addrList := flag.String("a", "", "address list (format ip:port ip:port)")
	timeout := flag.String("t", "3s", "connect timeout(default 3s)")
	flag.Parse()

	dialTimeout, err := time.ParseDuration(*timeout)
	if err != nil {
		panic(err)
	}

	addrs := make([]string, 0)
	if *addrListFileName != "" {
		b, err := ioutil.ReadFile(*addrListFileName)
		if err != nil {
			panic(err)
		}

		lines := strings.Split(string(b), "\n")
		addrs = addAddr(lines, addrs)
	}

	if *addrList != "" {
		lines := strings.Fields(*addrList)
		addrs = addAddr(lines, addrs)
	}

	okChan := make(chan string)
	errChan := make(chan Result)

	for _, addr := range addrs {
		go tcpCheck(addr, okChan, errChan, dialTimeout)
	}

	okAddrs := make([]string, 0)
	errAddrs := make([]Result, 0)

	cnt := len(addrs)
	for cnt > 0 {
		select {
		case addr := <-okChan:
			okAddrs = append(okAddrs, addr)
			cnt--
		case err := <-errChan:
			errAddrs = append(errAddrs, err)
			cnt--
		}
	}

	total := len(addrs)

	if len(errAddrs) > 0 {
		fmt.Printf("Total failed %d/%d :\n", len(errAddrs), total)
		for _, addr := range errAddrs {
			fmt.Printf("%s, error %v\n", addr.Addr, addr.Err)
		}

		fmt.Printf("\n\n")
	}

	if len(okAddrs) > 0 {
		fmt.Printf("Total ok %d/%d:\n", len(okAddrs), total)
		for _, addr := range okAddrs {
			fmt.Printf("%s\n", addr)
		}
	}
}

func addAddr(lines []string, addrs []string) []string {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if contains(addrs, line) {
			fmt.Printf("duplicate address %s\n", line)
		} else {
			addrs = append(addrs, line)
		}
	}
	return addrs
}

func tcpCheck(addr string, okChan chan string, errChan chan Result, timeout time.Duration) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		okChan <- addr
		_ = conn.Close()
	} else {
		errChan <- Result{addr, err}
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
