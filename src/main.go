package main

import (
	_ "go.uber.org/automaxprocs"
	"whoops/kafka2es/src/cmd"
)

func main() {
	//方法入口
	cmd.Run()
}
