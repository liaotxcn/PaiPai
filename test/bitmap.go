package main

import "fmt"

// 位运算
func main() {
	a := 7
	fmt.Printf("%b,%v \n", a, a)
	b := a << 2
	fmt.Printf("%b,%v \n", b, b)
	c := a >> 2
	fmt.Printf("%b,%v \n", c, c)

	d := a & b // 位与(只要有0结果就是0)
	fmt.Printf("%b,%v \n", d, d)
	e := a | b // 位或(只要有1结果就是1)
	fmt.Printf("%b,%v \n", e, e)
}
