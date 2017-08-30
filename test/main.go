package main
//
//import (
//	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
//	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
//	//"os"
//	"fmt"
//	"github.com/alecthomas/kingpin"
//
//	"strings"
//	"io/ioutil"
//	"os"
//	"bytes"
//)
//
//func main() {
//	// fmt.Println(os.Args)
//	vs := kingpin.Flag("v", "").Strings()
//	kingpin.Parse()
//	fmt.Println(*vs)
//	for _, v := range *vs {
//		//fmt.Printf("==>%s\n", v)
//		pairs := strings.Split(v, ":")
//		configMapName := pairs[0];
//		filePath := pairs[1];
//		fmt.Printf("==>%s\n", configMapName)
//		//fmt.Printf("%s--%s\n", configMapName, filePath)
//		if (!PathExists(filePath)) {
//			err := os.MkdirAll(filePath, 0777)
//			if err != nil {
//				panic(err.Error())
//			}
//		};
//		newData := []byte("xxx")
//		if (PathExists("/etc/nginx/conf-site.d/test.conf")) {
//			data, err := ioutil.ReadFile("/etc/nginx/conf-site.d/test.conf")
//			if err != nil {
//				panic(err.Error())
//			}
//			if (!bytes.Equal(data, newData)) {
//				err := ioutil.WriteFile("/etc/nginx/conf-site.d/test.conf", newData, 0644)
//				if err != nil {
//					panic(err.Error())
//				}
//			}
//		} else {
//			err := ioutil.WriteFile("/etc/nginx/conf-site.d/test.conf", newData, 0644)
//			if err != nil {
//				panic(err.Error())
//			}
//		}
//
//	}
//}
//
//func PathExists(path string) (bool) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true
//	}
//	if os.IsNotExist(err) {
//		return false
//	}
//	return false
//}
