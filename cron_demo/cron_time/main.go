package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	// 每隔5分钟执行一次
	// cronexpr支持7位数 秒 分 时 天 月 星期 年
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		// 每5秒
		fmt.Println(err)
		return
	}
	// 当前时间
	now = time.Now()
	// 下次调度时间
	nextTime = expr.Next(now)

	// 等待这个定时器超时(5秒后执行一个函数)
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了:", nextTime)
	})

	time.Sleep(5 * time.Second)
}
