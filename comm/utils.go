package comm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

//几个时间相关的函数
//获取当前时间字符串
func GetTimeString() string {
	return GetTimeStringByUTC(time.Now().Unix())
}

//获取utc时间字符串
func GetTimeStringByUTC(utc int64) string {
	t := time.Unix(utc, 0)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func GetTimeByString(tstr string) time.Time {
	if tstr <= "1970-01-01 00:00:00" {
		return time.Unix(0, 0)
	}
	//string转化为时间，layout必须为 "2006-01-02 15:04:05"
	t, err := time.Parse("2006-01-02 15:04:05", tstr)
	if err != nil {
		fmt.Printf("GetTimeByString err %v", err)
	}
	return t
}

//获取今日0点时间戳
func GetTodayBeginTime() int64 {
	t := time.Now()
	return t.Unix() - int64(t.Hour()*3600+t.Minute()*60+t.Second())
}

//获取年月,如:201703
func GetYearMonth() int {
	t := time.Now()
	s := fmt.Sprintf("%d%02d", t.Year(), t.Month())
	i, _ := strconv.Atoi(s)
	return i
}

func GetLogFlag() int {
	return log.LstdFlags | log.Llongfile
}

func GetMD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

//gameserver和loginserver通信签名用的key
const loginServerSignKey = "abcdefg"

func GetLoginServerSign(server_id int32) string {
	return GetMD5(loginServerSignKey + strconv.Itoa(int(server_id)))
}
