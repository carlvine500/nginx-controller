package nginx

import (
	"os"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api/v1"
	"io/ioutil"
	"bytes"
	"strings"
	"os/exec"
	"time"
	"math/rand"
)

func SyncConfigMapToLocalDir(clientset *kubernetes.Clientset, configmap2local *string) {
	glog.Infof("--configmap2local=%s", *configmap2local)
	configmap2locals := strings.Split(*configmap2local, ",")
	for _, v := range configmap2locals {
		pairs := strings.Split(v, ":")
		configMapName := pairs[0];
		localDir := pairs[1];
		if !PathExists(localDir) {
			err := os.MkdirAll(localDir, 0777)
			if err != nil {
				panic(err.Error())
			}
		}
		go watchConfigMap(clientset, configMapName, localDir);
	}
}

func watchConfigMap(clientset *kubernetes.Clientset, configMapName string, localDir string) {
	watcher, whErr := clientset.CoreV1().ConfigMaps("default").Watch(metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("metadata.name", configMapName).String()})
	if whErr != nil {
		print(whErr)
	}
	glog.Infof("watch configMap=%s,localDir=%s", configMapName, localDir)
	c := watcher.ResultChan()
	for {
		select {
		case e := <-c:
			v := reflect.ValueOf(e.Object)
			configMap, _ := v.Elem().Interface().(v1.ConfigMap)
			syncFile(configMap, localDir)
		}
	}
}

func syncFile(configMap v1.ConfigMap, localDir string) {
	hostname, _ := os.Hostname()
	canNginxReload := false
	localFileList, err := ioutil.ReadDir(localDir)
	if err != nil {
		glog.Errorf("readDir fail, localDir=%s,err=%v", localDir, err)
	}
	for _, fileInfo := range localFileList {
		if _, localFileExists := configMap.Data[fileInfo.Name()]; !localFileExists {
			localFilePath := localDir + "/" + fileInfo.Name();
			glog.Infof("remove localFilePath =%s", localFilePath)
			os.Remove(localFilePath)
			// canNginxReload = true
		}
	}
	for fileName, fileContent := range configMap.Data {
		localFilePath := localDir + "/" + fileName;
		newData := []byte( fileContent)

		if strings.Compare(fileName, "hostname_date") == 0 {
			if strings.Contains(fileContent, hostname) {
				canNginxReload = true
			}
			oldData, _ := ioutil.ReadFile(localFilePath)
			if !bytes.Equal(oldData, newData) {
				canNginxReload = true
			}
		}

		if !PathExists(localFilePath) {
			err := ioutil.WriteFile(localFilePath, newData, 0644)
			if err != nil {
				glog.Errorf("fist time write fail, localFilePath =%s,err=%v", localFilePath, err)
			} else {
				glog.Infof("fist time write localFilePath =%s", localFilePath)
			}
		} else {
			oldData, err := ioutil.ReadFile(localFilePath)
			if err != nil {
				glog.Errorf("read fail, localFilePath=%s,err=%v", localFilePath, err)
			}
			if !bytes.Equal(oldData, newData) {
				err := ioutil.WriteFile(localFilePath, newData, 0644)
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

	if canNginxReload {
		reloadNginx()
	}
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

func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
