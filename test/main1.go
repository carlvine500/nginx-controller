package main

import (
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	//"os"

	"fmt"
	"os/exec"
	//"strings"
	//"io/ioutil"
)

func main2() {
	//cmd := exec.Command("/bin/sh","-c","nginx -s reload")
	cmd := exec.Command("/bin/sh","-c","nginx -t")
	re, err := cmd.CombinedOutput()
	//stdout, err := cmd.StdoutPipe()
	//bytes, err := ioutil.ReadAll(stdout)
	//fmt.Printf("result=%v\n",bytes)
	//if err != nil {
	//	fmt.Println("ReadAll stdout: ", err.Error())
	//	//return
	//}
	//successful
	//if strings.Contains(string(re),"CONTAINER"){
	//	fmt.Printf("CONTAINER\n")
	//}
	fmt.Printf("result=%s\n",re)
	fmt.Printf("err=%v\n",err)

}
