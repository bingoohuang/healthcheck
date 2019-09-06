# healthcheck
healthcheck cli for tcp, http and etc.

## Install

`go get github.com/bingoohuang/healthcheck/...`

## Usage

```bash
➜  Downloads healthcheck -h
Usage of healthcheck:
  -a string
    	address list (format ip:port ip:port)
  -f string
    	address file name(ip:port per line)
  -t string
    	connect timeout(default 3s) (default "3s")
    	
    	
➜  Downloads healthcheck -a 123.206.185.162:53306 -f ip-port.txt -t 1s
Total failed 34/35 :
192.168.20.5:10983, error dial tcp 192.168.20.5:10983: i/o timeout
192.168.29.17:10017, error dial tcp 192.168.29.17:10017: i/o timeout
192.168.22.25:10616, error dial tcp 192.168.22.25:10616: i/o timeout
192.168.20.5:18080, error dial tcp 192.168.20.5:18080: i/o timeout
192.168.22.6:9082, error dial tcp 192.168.22.6:9082: i/o timeout
```