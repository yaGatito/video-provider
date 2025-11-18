package main

import (
	"fmt"
)

func main() {
	n := 2149583361
	fmt.Println(Int32ToIp(uint32(n)))
}

func Int32ToIp(n uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", (n>>24)&0b11111111, (n>>16)&0b11111111, (n>>8)&0b11111111, n&0b11111111)
}

