package nginx

import (
	"bytes"
	"context"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

func init() {
	glog.MaxSize = 1024 * 1024 * 200
}

func SyncConfigMapToLocalDir(clientset *kubernetes.Clientset, configmap2local *string) {
	glog.Infof("--configmap2local=%s", *configmap2local)
	configmap2locals := strings.Split(*configmap2local, ",")
	for _, v := range configmap2locals {
		pairs := strings.Split(v, ":")
		configMapName := pairs[0]
		localDir := pairs[1]
		if !PathExists(localDir) {
			err := os.MkdirAll(localDir, 0777)
			if err != nil {
				panic(any(err.Error()))
			}
		}
		go watchConfigMap2(clientset, configMapName, localDir)
	}
}

func watchConfigMap2(clientset *kubernetes.Clientset, configMapName string, localDir string) {
	//for {
		watchConfigMap(clientset, configMapName, localDir)
	//}
}

func watchConfigMap(clientset *kubernetes.Clientset, configMapName string, localDir string) {
	defer func() {
		if err := recover(); any(err) != nil {
			glog.Errorf("Unknow Error[E],err=%v,configMapName=%s,localDir=%s", err, configMapName, localDir)
		}
	}()
	options := metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("metadata.name", configMapName).String()}
	watcher, whErr := clientset.CoreV1().ConfigMaps("default").Watch(context.TODO(), options)
	if whErr != nil {
		print(whErr)
		return
	}
	glog.Infof("watch configMap=%s,localDir=%s", configMapName, localDir)

	for event := range watcher.ResultChan() {
		if event.Type == watch.Modified || event.Type == watch.Added {
			configMap, ok := event.Object.(*v1.ConfigMap)
			if !ok {
				glog.Errorf("Error: Could not decode object,configMapName=%s", configMapName)
				continue
			}
			glog.Infof("ConfigMap=%s %s\n", configMap.Name, event.Type)
			// eg: /etc/nginx/conf-ssl.d/a.txt==> xy
			// configMap.Data["a.txt"]=xy
			syncFile(*configMap, localDir)
		}
	}
}

func syncFile(configMap v1.ConfigMap, localDir string) {
	hostname, _ := os.Hostname()
	//hostname:= "ip-172-37-100-93"
	//canNginxReload := false
	localFileList, err := os.ReadDir(localDir)
	if err != nil {
		glog.Errorf("readDir fail, localDir=%s,err=%v", localDir, err)
	}

	//ignore override source host nginx config
	valueOfHostnameDate,fileExists:= configMap.Data["hostname_date"];
	if !fileExists {
		return
	}
	valueOfHostnameDate = strings.TrimSpace(valueOfHostnameDate)

	contentOfHostnameDate, _ := os.ReadFile(localDir+"/"+"hostname_date")
	oldValueOfHostnameDate :=strings.TrimSpace(string(contentOfHostnameDate))
	//hostname who upload configMap reload in n seconds
	if strings.Contains(valueOfHostnameDate, hostname) {
		timeString := valueOfHostnameDate[strings.Index(valueOfHostnameDate, " ")+1:]
		if subSecondsFromNow(timeString) < 30 {
			reloadNginx()
		}
		return
	}

	//hostname who receive configMap
	if strings.Compare(valueOfHostnameDate, oldValueOfHostnameDate) == 0 {
		//canNginxReload = true
		return
	}

	tmpDir := localDir + "/tmp"
	_, err2 := os.Stat(tmpDir)
	if os.IsNotExist(err2) {
		os.Mkdir(tmpDir, os.ModePerm)
	}

	// don't delete config,mv it to tmp/ directory
	for _, fileInfo := range localFileList {
		if _, localFileExists := configMap.Data[fileInfo.Name()]; !localFileExists {
			localFilePath := localDir + "/" + fileInfo.Name()
			tmpFilePath := tmpDir + "/" + fileInfo.Name()
			glog.Infof("mv localFilePath to %s", tmpFilePath)
			// os.Remove(localFilePath)
			os.Rename(localFilePath, tmpFilePath)
			// canNginxReload = true
		}
	}

	// override local config from configMap
	for fileName, fileContent := range configMap.Data {
		localFilePath := localDir + "/" + fileName
		newData := []byte(fileContent)

		if !PathExists(localFilePath) {
			err := os.WriteFile(localFilePath, newData, 0644)
			if err != nil {
				glog.Errorf("fist time write fail, localFilePath =%s,err=%v", localFilePath, err)
			} else {
				glog.Infof("fist time write localFilePath =%s", localFilePath)
			}
		} else {
			oldData, err := os.ReadFile(localFilePath)
			if err != nil {
				glog.Errorf("read fail, localFilePath=%s,err=%v", localFilePath, err)
			}
			if !bytes.Equal(oldData, newData) {
				err := os.WriteFile(localFilePath, newData, 0644)
				if err != nil {
					glog.Errorf("configMap file changed,but write fail, localFilePath=%s,err=%v", localFilePath, err)
				} else {
					glog.Infof("configMap file changed,write to localFilePath=%s", localFilePath)
				}
			} else {
				glog.Infof("configMap file is the same as localFilePath=%s", localFilePath)
			}
		}
	}
	//if canNginxReload {
	reloadNginx()
	//}
}

func reloadNginx() {
	// sleep random seconds to avoid global performance effect
	time.Sleep(1 + time.Duration(rand.Intn(3))*time.Second)
	test := exec.Command("/bin/sh", "-c", "nginx -t")
	testResult, testErr := test.CombinedOutput()
	if testErr != nil {
		glog.Errorf("nginx -t err=%v", testErr)
	}
	glog.Infof("nginx -t result=%s", testResult)
	if !strings.Contains(string(testResult), "successful") {
		return
	}
	reload := exec.Command("/bin/sh", "-c", "nginx -s reload")
	reloadResult, err := reload.CombinedOutput()
	glog.Infof("nginx -s reload result=%s", reloadResult)
	if err != nil {
		glog.Errorf("nginx -s reload err=%v", err)
	}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}


func subSecondsFromNow(unixTimeString string) (float64) {
	unixTime, err := time.Parse(time.UnixDate, unixTimeString)
	//glog.Errorf("unixTimeString=%s", unixTimeString)
	if err != nil {
		glog.Errorf("error time string=%s", unixTimeString)
		return 0
	}
	now := time.Now()
	glog.Infof("now time=%v",now)
	glog.Infof("file time=%v",unixTime)
	minutes := now.Sub(unixTime).Seconds()
	return minutes
}

