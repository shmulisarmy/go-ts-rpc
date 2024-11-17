package main

import (
	rpc "github.com/shmulisarmy/go-ts-rpc"
)

func add(a int, b int) int {
	return a + b
}

func main() {
	rpc.Load_file("main.go")
	rpc.Setup_rpc(add)
}
