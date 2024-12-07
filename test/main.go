package main

import (
	"golang.org/x/exp/errors/fmt"
	"strings"
	"time"
)

func main1() {
	//re,_:=ioutil.ReadFile("/tmp/abc")
	//fmt.Print(bytes.Equal([]byte("hostname"),re))
	//for i := 0; i < 10; i++ {
	//	t := 3 + rand.Intn(3)
	//	time.Sleep(time.Duration(t) * time.Second)
	//	fmt.Println(t)
	//
	//
	//}
	//_, err := os.Stat("/Users/tingfeng/work/golang/tmp")
	//if os.IsNotExist(err) {
	//	os.Mkdir("/Users/tingfeng/work/golang/tmp", os.ModePerm)
	//}
	//err1 := os.Rename("/Users/tingfeng/work/golang/a", "/Users/tingfeng/work/golang/tmp/a")
	//fmt.Print(err1)
	//flysnowRegexp := regexp.MustCompile(`\d+:\d+:\d+`)
	s := "ip-172-37-100-93 Tue Dec 1 22:40:57 CST 2020\n"
	s = strings.TrimSpace(s)
	timeString := s[strings.Index(s, " ")+1 : len(s)]
	m := subMinutesFromNow(timeString)
	fmt.Println(timeString)
	fmt.Println(m)
	//params := flysnowRegexp.FindStringSubmatch(s)
	//a := subMinutesFromNow(params[0])
	//fmt.Println(a)

}

func subMinutesFromNow(unixTimeString string) (float64) {
	unixTime, err := time.Parse(time.UnixDate, unixTimeString)
	if err != nil {
		return 0
	}
	now := time.Now()
	fmt.Println("now time=",now)
	fmt.Println("file time=",unixTime)
	minutes := now.Sub(unixTime).Minutes()
	return minutes
}
