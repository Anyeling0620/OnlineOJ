package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	// go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("stdinPipe err")
		panic(err)
	}
	_, err = io.WriteString(stdinPipe, "11 23\n")
	if err != nil {
		fmt.Println("io.WriteString err")
		panic(err)
	}

	if err = stdinPipe.Close(); err != nil {
		fmt.Println("close stdinPipe err")
	}
	// 根据测试的输入案例进行运行，拿到输出结果和标准输出结果进行比对
	if err := cmd.Run(); err != nil {
		fmt.Println("stderr:", stderr.String())
		fmt.Println("stdout:", out.String())
	}
	fmt.Println(out.String())

	println(out.String() == "34\n")
}
