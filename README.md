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

## 定期探测网络是否成功

1. 构建 `go build -ldflags="-s -w" ./cmd/healthcheck`
1. 建立.env文件

    ```properties
    # 应用名称，默认使用当前pid
    APP_NAME=healthcheck
    # 写入指标日志的间隔时间，默认1s
    METRICS_INTERVAL=1s
    # 写入心跳日志的间隔时间，默认20s
    HB_INTERVAL=20s
    # Metrics对象的处理容量，默认1000，来不及处理时，超额扔弃处理
    CHAN_CAP=1000
    # 日志存放的目录，默认/tmp/log/metrics
    LOG_PATH=/var/log/footstone/metrics
    # 日志文件最大保留天数
    MAX_BACKUPS=7
    ```

1. 建立地址列表文件addr.txt

    ```
    192.168.126.16:22 测试机1
    192.168.126.18:22 测试机2
    ```

1. 运行程序

    `healthcheck -interval=10s -f addr.txt`

1. 查看输出

    .../metrics/metrics-hb.healthcheck.log:
    
    ```json
    {"time":"20201111135251000","key":"healthcheck.hb","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"HB","v1":1,"v2":0,"min":0,"max":0}
    {"time":"20201111135311000","key":"healthcheck.hb","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"HB","v1":1,"v2":0,"min":0,"max":0}
    ```
    
    .../metrics/metrics-key.healthcheck.log:
    
    ```json
    {"time":"20201111143325000","key":"测试机1#192_168_126_16#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"SUCCESS_RATE","v1":1,"v2":1,"min":100,"max":100}
    {"time":"20201111143325000","key":"HealthCheck#192_168_126_182#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"QPS","v1":1,"v2":0,"min":-1,"max":-1}
    {"time":"20201111143325000","key":"HealthCheck#192_168_126_182#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"FAIL_RATE","v1":0,"v2":1,"min":-1,"max":-1}
    {"time":"20201111143325000","key":"测试机2#192_168_126_18#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"FAIL_RATE","v1":0,"v2":1,"min":-1,"max":-1}
    {"time":"20201111143325000","key":"HealthCheck#192_168_126_182#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"SUCCESS_RATE","v1":1,"v2":1,"min":100,"max":100}
    {"time":"20201111143325000","key":"测试机1#192_168_126_16#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"QPS","v1":1,"v2":0,"min":-1,"max":-1}
    {"time":"20201111143325000","key":"测试机1#192_168_126_16#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"FAIL_RATE","v1":0,"v2":1,"min":-1,"max":-1}
    {"time":"20201111143325000","key":"测试机2#192_168_126_18#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"SUCCESS_RATE","v1":1,"v2":1,"min":100,"max":100}
    {"time":"20201111143325000","key":"测试机2#192_168_126_18#22","hostname":"bingoobjcadeMacBook-Pro.local","logtype":"QPS","v1":1,"v2":0,"min":-1,"max":-1}
    ```

## Resouces

1. [areyouok](https://github.com/Bhupesh-V/areyouok) A fast and easy to use URL health checker ⛑️ Keep your links healthy during tough times (Out of box support for GitHub Actions)
