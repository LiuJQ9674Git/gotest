package main
//go语言处理的是包,而不是文件.所有的go语言必须在一个包中
//包名和函数名不会有命名冲突
import (
	"fmt"//引入标准包
)

func show() (err error) {
	//err="tt"
	return nil
}

func main() {
	fmt.Println("Hello World!")

}