package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)
	//cmd = exec.Command("/bin/bash", "-c", "echo 1;echo 2;")
	cmd = exec.Command("D:\\Cygwin\\bin\\bash.exe", "-c", "echo 1")
	err = cmd.Run()
	fmt.Println(err)
}
