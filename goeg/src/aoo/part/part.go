package main

import (
	"strings"
	"fmt"
)
//结构内部定义为指针，不是在方法上而是直接使用结构的变量，
//则符合GO关于值、引用和指针三种类型的定义
type PartIDPoint struct {
	Id *int //直接修改指针变量，一改全改，类似于全局变量,但是SetID时没有修改
	Name string
	age int
}

//需要修改Part结构的值，因此需要接收者为指针
func (part *PartIDPoint) SetId(id *int) {
	part.Id=id
}

//需要修改Part结构的值，因此需要接收者为指针
func (part PartIDPoint) SetIdValue(id *int) {
	part.Id=id
}

//需要修改Part结构的值，因此需要接收者为指针
func (part *PartIDPoint) SetAge(age int) {
	part.age=age
}

//需要修改Part结构的值，因此需要接收者为指针
func (part *PartIDPoint) UpperCase() {
	part.Name=strings.ToUpper(part.Name)
}

//需要修改Part结构的值，因此需要接收者为指针,如果接收者为值则不能修改
func (part PartIDPoint) LowerCase() {
	part.Name=strings.ToLower(part.Name)
}

func (part PartIDPoint) String() string{
	//因为变量为指针，所以需要获取指针值的
	return fmt.Sprintf("%d %q %d",*part.Id,part.Name,part.age)
}

//结构内部定义为指针，不是在方法上而是直接使用结构的变量，

type PartIDValue struct {
	Id int
	Name string
}

//需要修改Part结构的值，因此需要接收者为指针
func (part *PartIDValue) UpperCase() {
	part.Name=strings.ToUpper(part.Name)
}

//需要修改Part结构的值，因此需要接收者为指针,如果接收者为值则不能修改
func (part PartIDValue) LowerCase() {
	part.Name=strings.ToLower(part.Name)
}

func (part PartIDValue) String() string{
	//因为变量为指针，所以需要获取指针值的
	return fmt.Sprintf("%d %q",part.Id,part.Name)
}

func main()  {
	id:=15
	//值
	part:=PartIDPoint{&id,"test",90}
	fmt.Println("值测试开始")
	//另一个结构的指针
	partSecond:=&PartIDPoint{&id,"partSecond",10}
	id=1

	fmt.Println("part is->",part)
	id=10
	part.Id=&id
	part.UpperCase()
	fmt.Println("UpperCase part is->",part)
	//由于接收者为值所以没有修改结构内的变量Name
	part.LowerCase()
	id=20
	part.Id=&id //因为Id定义为指针，所以
	fmt.Println("LowerCase part is->",part)
	//值对象
	part.SetAge(30)
	fmt.Println("SetAge->",part)
	//修改id
	var vId=999;
	var pId=&vId
	part.SetId(pId)
	fmt.Println("SetId->",part)
	vId=1000
	part.SetIdValue(pId)
	fmt.Println("SetId->",part)
	//另一个结构的值
	//partSecond:=Part{15,"partSecond"}
	//由于ID为指针类型，所以partSecond的ID为最后设置的20
	fmt.Println("指针测试开始")
	fmt.Println("part is->",partSecond)

	//值对象
	partSecond.SetAge(30)
	fmt.Println("SetAge->",partSecond)
	//如果修改PartIDValue，则part
	// Second的id值为初始化设置值
	//下面更符合面向对象的定义
	id=15
	//另一个结构的值
	partSecondValue:=PartIDValue{id,"partSecond"}
	id=1
	partValue:=PartIDValue{id,"test"}
	fmt.Println("part is->",partValue)
	id=10
	partValue.Id=id
	partValue.UpperCase()
	fmt.Println("UpperCase part is->",partValue)
	//由于接收者为值所以没有修改结构内的变量Name
	partValue.LowerCase()
	id=20
	partValue.Id=id
	fmt.Println("LowerCase part is->",partValue)
	//另一个结构的值
	//partSecond:=Part{15,"partSecond"}
	//由于ID为指针类型，所以partSecond的ID为最后设置的20
	fmt.Println("part is->",partSecondValue)

	//如果修改PartIDValue，则partSecond的id值为初始化设置值
}