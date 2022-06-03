package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	// 生成Cmd
	cmd = exec.Command("D:\\Cygwin\\bin\\bash.exe", "-c", "ls")

	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		fmt.Println("error")
		return
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}
