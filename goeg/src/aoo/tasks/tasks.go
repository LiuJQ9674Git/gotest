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

type Task struct {
	slice []string
	Count //嵌入,当定义为*Count编译报错，报错如下
	//panic: runtime error: invalid memory address or nil pointer dereference
}

func (tasks *Task) Add(t string) {
	tasks.slice=append(tasks.slice,t)
	tasks.Increment()

}

func (tasks *Task) Pop() {
	tasks.slice=tasks.slice[1:]
	tasks.Decrement()
}


func main()  {
	task :=Task{}
	task.Add("test")
	fmt.Println(task)

	task.Add("test2")
	fmt.Println(task)

	task.Pop()
	fmt.Println(task)

	task.Pop()
	fmt.Println(task)
}