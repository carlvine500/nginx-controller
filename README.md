# nginx-controller
sychromize nginx config file by kubernetes.configMap,and reload nginx when file(hostname_date) changed
- step1: upload nginx config files to configMap by kubectl
- step2: nginx-controller will listen to configMap, and download nginx config files to local directory
- step3: test and reload nginx when files changed
- note: the file "hostname_date" in any directory, used for force reload nginx

# build
- linux `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build`
- macos `go build`

# example
```
$cd /etc/nginx
$tree
├── conf-site.d
│   └── test.conf
├── conf-ssl.d
│   └── ssl_session_ticket.key
├── conf-upstream.d
│   ├── a.conf
│   └── b.conf
├── conf-socks.d
│   ├── a.conf
│   └── b.conf
└── publish.sh

$cat publish.sh #copy files to configMap,update file hostname_date
echo `hostname` `date` | sudo  tee /etc/nginx/conf-site.d/hostname_date
kubectl create configmap nginx-ssl --from-file /etc/nginx/conf-ssl.d -o yaml --dry-run | kubectl apply -f -
kubectl create configmap nginx-upstream --from-file /etc/nginx/conf-upstream.d -o yaml --dry-run | kubectl apply -f -
kubectl create configmap nginx-site --from-file /etc/nginx/conf-site.d -o yaml --dry-run | kubectl apply -f -
kubectl create configmap nginx-socks --from-file /etc/nginx/conf-socks.d -o yaml --dry-run | kubectl apply -f -
tail -f -n 1 /var/log/nginx/nginx-controller.INFO

$kubectl get configmap 
nginx-site       2         1d
nginx-ssl        1         28s
nginx-upstream   2         14h

$./nginx-controller #download file to local directory,reload nginx when hostname_date updated
watch configMap=nginx-ssl,directory=/etc/nginx/conf-ssl.d 
watch configMap=nginx-upstream,directory=/etc/nginx/conf-upstream.d 
watch configMap=nginx-site,directory=/etc/nginx/conf-site.d 
```

# attention
- suggest file encoding=base64 , becauseof kubernetes.client-go ConfigMap's value is string 
- `openssl rand -base64 48 > ssl_session_ticket.key`

# Usage of ./nginx-controller:
```
  -alsologtostderr
        log to standard error as well as files
  -configmap2local string
        configMap:localDir, eg:nginx-site:/etc/nginx/conf-site.d,nginx-upstream:/etc/nginx/conf-upstream.d,nginx-ssl:/etc/nginx/conf-ssl.d (default "nginx-site:/etc/nginx/conf-site.d,nginx-upstream:/etc/nginx/conf-upstream.d,nginx-ssl:/etc/nginx/conf-ssl.d")
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "/Users/tingfeng/.kube/config")
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
