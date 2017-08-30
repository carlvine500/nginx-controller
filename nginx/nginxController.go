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
)

func SyncConfigMapToDirectory(clientset *kubernetes.Clientset, configmap2file *string) {
	glog.Infof("--configmap2file=", *configmap2file)
	configmap2files := strings.Split(*configmap2file, ",")
	for _, v := range configmap2files {
		pairs := strings.Split(v, ":")
		configMapName := pairs[0];
		directory := pairs[1];
		if (!PathExists(directory)) {
			err := os.MkdirAll(directory, 0777)
			if err != nil {
				panic(err.Error())
			}
		}
		go watchConfigMap(clientset, configMapName, directory);
	}
}

func watchConfigMap(clientset *kubernetes.Clientset, configMapName string, directory string) {
	watcher, whErr := clientset.CoreV1().ConfigMaps("default").Watch(metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("metadata.name", configMapName).String()})
	if whErr != nil {
		print(whErr)
	}
	glog.Infof("watch configMap=%s,directory=%s", configMapName, directory)
	c := watcher.ResultChan()
	for {
		select {
		case e := <-c:
			v := reflect.ValueOf(e.Object)
			configMap, _ := v.Elem().Interface().(v1.ConfigMap)
			syncFile(configMap, directory)
		}
	}
}

func syncFile(configMap v1.ConfigMap, directory string) {
	canNginxReload := false
	fileList, err := ioutil.ReadDir(directory)
	if err != nil {
		glog.Errorf("readDir fail, directory=%s,err=%v", directory, err)
	}
	for _, fileInfo := range fileList {
		value := configMap.Data[fileInfo.Name()]
		realFilePath := directory + "/" + fileInfo.Name();
		if len(value) == 0 {
			glog.Infof("remove realFilePath =%s", realFilePath)
			os.Remove(realFilePath)
			canNginxReload = true
		}
	}
	for fileName, fileContent := range configMap.Data {
		realFilePath := directory + "/" + fileName;
		newData := []byte( fileContent)
		if !PathExists(realFilePath) {
			err := ioutil.WriteFile(realFilePath, newData, 0644)
			if err != nil {
				glog.Errorf("fist time write fail, realFilePath =%s,err=%v", realFilePath, err)
			} else {
				glog.Infof("fist time write realFilePath =%s", realFilePath)
				canNginxReload = true
			}
		} else {
			oldData, err := ioutil.ReadFile(realFilePath)
			if err != nil {
				glog.Errorf("read fail, realFilePath=%s,err=%v", realFilePath, err)
			}
			if !bytes.Equal(oldData, newData) {
				err := ioutil.WriteFile(realFilePath, newData, 0644)
				if err != nil {
					glog.Errorf("configMap file changed,but write fail, realFilePath=%s,err=%v", realFilePath, err)
				} else {
					glog.Infof("configMap file changed,write to realFilePath=%s", realFilePath)
					canNginxReload = true
				}
			} else {
				glog.Infof("configMap file is the same as realFilePath=%s", realFilePath)
			}
		}

		if canNginxReload {
			cmd := exec.Command("nginx", "-s", "reload")
			_, err := cmd.Output()
			if err != nil {
				glog.Errorf("nginx reload err=%v", err)
			}
		}
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
