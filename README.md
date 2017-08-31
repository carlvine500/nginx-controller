# nginx-controller

# build
- linux `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build`
- macos `go build`

# kubectl
`cd /etc/nginx/conf-site;./publish.sh`

```
Usage of ./nginx-controller:
  -alsologtostderr
        log to standard error as well as files
  -configmap2file string
        configMap:directory, eg:nginx-site:/etc/nginx/conf-site.d,nginx-upstream:/etc/nginx/conf-upstream.d (default "nginx-site:/etc/nginx/conf-site.d,nginx-upstream:/etc/nginx/conf-upstream.d")
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "/root/.kube/config")
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```
