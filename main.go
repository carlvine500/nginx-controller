package main

import (
	"flag"
	"github.com/carlvine500/nginx-controller/nginx"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	defaultConfig :=
		"nginx-site:/etc/nginx/conf-site.d" +
			",nginx-upstream:/etc/nginx/conf-upstream.d" +
			",nginx-ssl:/etc/nginx/conf-ssl.d"+
			",nginx-socks:/etc/nginx/conf-socks.d"
	configmap2local := flag.String("configmap2local", defaultConfig, "configMap:localDir, eg:"+defaultConfig)

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(any(err.Error()))
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(any(err.Error()))
	}
	//re, _ := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	//fmt.Printf("namespaces=%v",re)

	nginx.SyncConfigMapToLocalDir(clientset, configmap2local)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
