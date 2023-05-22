package main

import "fmt"

type Count int

func (count *Count)  Increment(){
	*count++
}
//方法签名如果修改为值为接收者，则Decrement执行完之后，count值不变
// func (count Count)  Decrement()
func (count *Count)  Decrement(){
	*count--
}

func (count Count)  IsZero() bool{
	return count==0
}

func main()  {
	//Count是可寻址的，指针接收的方法也可以用在值上，此时GO会解析为指针。
	//可寻址：解引用指针、变量、切片、结构的变量

	count:=Count(2);
	fmt.Println(count)
	zero:=count.IsZero()
	fmt.Println(zero)
	count.Decrement()
	fmt.Println(count)
	dezero:=count.IsZero()
	fmt.Println(dezero)

	count.Decrement()
	fmt.Println(count)
	dezero=count.IsZero()
	fmt.Println(dezero)
}