package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 设置系统时间
func setSystemTime(t time.Time) error {
	// 格式化时间字符串
	timeStr := t.Format("20060102 15:04:05.00")

	// 拆分时间
	parts := strings.Split(timeStr, " ")
	dateStr, timeStr := parts[0], parts[1]

	// 拆分日期
	parts = strings.Split(dateStr, "")
	year, month, day := parts[0]+parts[1]+parts[2]+parts[3], parts[4]+parts[5], parts[6]+parts[7]

	// 拆分时间
	parts = strings.Split(timeStr, "")
	hour, minute, second := parts[0]+parts[1], parts[3]+parts[4], parts[6]+parts[7]

	// 构造命令行参数并执行
	cmd := exec.Command("cmd", "/C", "date", year+"/"+month+"/"+day)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("cmd", "/C", "time", hour+":"+minute+":"+second)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// 主函数
func main() {
	// 设置NTP服务器地址
	server := "ntp.aliyun.com"

	for {
		// 获取当前本地时间
		localTime := time.Now()

		// 获取NTP时间
		ntpTime, err := ntp.Time(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error syncing system time with NTP server %s: %s\n", server, err)
			time.Sleep(time.Minute)
			continue
		}

		// 计算时间差并设置本地时间
		delta := ntpTime.Sub(localTime)
		err = setSystemTime(localTime.Add(delta))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting system time: %s\n", err)
			time.Sleep(time.Minute)
			continue
		}

		fmt.Printf("System time adjusted by %v\n", delta)

		time.Sleep(time.Minute)
	}
}
