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
    	
    	
➜  Downloads healthcheck -a 123.206.185.162:53306 -f ~/ip-port.txt -t 1s
Failed 2/3:
192.168.37.107:10502, error dial tcp 192.168.37.107:10502: i/o timeout
192.168.20.5:18080, error dial tcp 192.168.20.5:18080: i/o timeout

OK 1/3:
123.206.185.162:53306
```