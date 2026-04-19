package core

import "fmt"

func Jc(x int) int { //定义一个递归函数
	var z int
	if x < 0 {
		fmt.Printf("你输入的数字有误")
	} else if x == 1 || x == 0 { //给予递归结束（判断 x 的值是否为1,0）---条件成了后将 x 的值 倒回去重新计算
		z = 1
	} else {
		z = Jc(x-1) * x //递归--调用自身函数
	}
	return z
}
