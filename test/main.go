package main

import (
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	//"os"

	//"strings"
	//"io/ioutil"
	"fmt"
	//"strings"
	"bytes"
	"io/ioutil"
	"os"
)

func main() {
	re,_:=ioutil.ReadFile("/tmp/abc")
	fmt.Print(bytes.Equal([]byte("hostname"),re))
	//for i := 0; i < 10; i++ {
	//	t := 3 + rand.Intn(3)
	//	time.Sleep(time.Duration(t) * time.Second)
	//	fmt.Println(t)
	//
	//
	//}
	_, err := os.Stat("/Users/tingfeng/work/golang/tmp")
	if os.IsNotExist(err) {
		os.Mkdir("/Users/tingfeng/work/golang/tmp", os.ModePerm)
	}
	err1 := os.Rename("/Users/tingfeng/work/golang/a", "/Users/tingfeng/work/golang/tmp/a")
	fmt.Print(err1)
}
