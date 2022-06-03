package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

// CronJob 代表一个任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	// 需要有1个调度协程, 它定时检查所有的Cron任务, 谁过期了就执行谁

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		err           error
		now           time.Time
		scheduleTable map[string]*CronJob // key: 任务的名字
	)

	scheduleTable = make(map[string]*CronJob)

	// 当前时间
	now = time.Now()

	// 定义两个cronjob
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	// 任务注册到调度表
	scheduleTable["job_1"] = cronJob

	// 定义两个cronjob
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	// 任务注册到调度表
	scheduleTable["job_2"] = cronJob

	// 启动调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
		)
		// 定时检查一下任务调度表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					// 计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, " 下次执行时间:", cronJob.nextTime)
				}
			}

			// 睡眠100毫秒
			select {
			case <-time.NewTimer(100 * time.Microsecond).C:
			}
		}
	}()

	time.Sleep(100 * time.Second)
}
