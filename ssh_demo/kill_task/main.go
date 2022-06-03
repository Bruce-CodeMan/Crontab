package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	// 执行1个cmd命令, 让他在一个协程里面去执行, 让他执行2秒, sleep 2; echo hello

	// 1秒的时候, 我们杀死cmd
	var (
		ctx           context.Context
		cancelFunc    context.CancelFunc
		cmd           *exec.Cmd
		resultChannel chan *result
		res           *result
	)

	resultChannel = make(chan *result, 1000)

	ctx, cancelFunc = context.WithCancel(context.TODO())

	// context: chan byte
	// 返回 cancelFunc : close(chan byte)

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "D:\\Cygwin\\bin\\bash.exe", "-c", "sleep 2;echo hello;")

		// ctx的内部维护了select { case <- ctx.Done():}
		// kill pid 进程ID

		// 执行任务, 捕获输出
		output, err = cmd.CombinedOutput()

		// 把任务输出结果传给main协程
		resultChannel <- &result{
			err:    err,
			output: output,
		}
	}()

	// 继续往下走
	time.Sleep(3 * time.Second)

	// 取消上下文
	cancelFunc()

	// 在main函数中等待子协程的退出, 并打印任务执行结果
	res = <-resultChannel

	fmt.Println(res.err, string(res.output))
}
