package main

import (
	"golang.org/x/exp/errors/fmt"
	"regexp"
	"time"
)

func main() {
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
	flysnowRegexp := regexp.MustCompile(`\d+:\d+:\d+`)
	params := flysnowRegexp.FindStringSubmatch("ip-172-37-100-93 Tue Dec 1 22:12:32 CST 2020")
	a := subMinutesFromNow(params[0])
	fmt.Println(a)

}

func subMinutesFromNow(hourMinuteSecond string) (float64) {
	now := time.Now()
	format := "2006-01-02 15:04:05"

	nowString := now.Format(format)
	yearMonthDay := now.Format("2006-01-02")

	nowUTC, _ := time.Parse(format, nowString)
	t2, _ := time.Parse(format, yearMonthDay+" "+hourMinuteSecond)

	minutes := nowUTC.Sub(t2).Minutes()
	return minutes
}
