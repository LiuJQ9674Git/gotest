package main

import "fmt"

//如果在一个值类型上调用某个方法，而该方法又需要一个值的指针参数，那么GO会自动的将该值
//的地址传递给方法。
//如果在某个值的指针上调用值的方法，那么该方法则实际需要一个值，此时GO会将此指针解引用，
//将该指针所指向的值传递给方法
type Exchanger interface {
	Exchange()
}

//切片结构
type Point [2]int

func (pair *Point) Exchange() {
	pair[0],pair[1]=pair[1],pair[0]
}

//自定义结构
type StringPair struct {
	first,second string
}

//Exchange方法接收者为指针，修改调用该方法的（指针所指向的）值
//jerryll.Exchange()签名为指针，即下面签名时,调用Exchange方法后值会变化
//(pair *StringPair) Exchange()
//如果签名如下，则调用前后值不变
//(pair *StringPair) Exchange()
func (pair *StringPair) Exchange() {
	pair.first,pair.second=pair.second,pair.first
}

//自定义结构，如果为指针，使用stringPoint.third设置时影响范围为全局，见part
type StringPoint struct {
	third *string
}

//SetThird的接收者为指针时，结构内变量值可以随着变化而变化。
//func (pair *StringPoint) SetThird(third *string)
//接收者为值时，结构内变量值不变，如下
//func (pair *StringPoint) SetThird(third *string)
func (pair *StringPoint) SetThird(third *string) {
	pair.third=third;
}

func (pair StringPoint) String() string{
	//因为变量为指针，所以需要获取指针值的
	return fmt.Sprintf("%q",*pair.third)
}

//自定义结构，结构内部的值是否为指针，GO可以根据可寻址来自动识别
//是否需要修改和接收者是指针类型还是值类型相关。
//在GO中，通道、切片、映射、函数和接口使用make定义，返回一个指向特定类型值的引用。
//定义语法为 new(Type)等价于&Type{}
type StringValue struct {
	third string
}

func (pair *StringValue) SetThird(third string) {
	pair.third=third;
}

func (pair StringValue) String() string{
	//因为变量为指针，所以需要获取指针值的
	return fmt.Sprintf("%q",pair.third)
}


func main() {
	fmt.Println("StringPair")
	//下面两种GO可以解析，结果等价
	//jerryll:=&StringPair{"Henery","Jerryll"}
	jerryll:=StringPair{"Henery","Jerryll"}
	fmt.Println(jerryll)
	// 虽然定义为值类型jerryll，方法的接收者为指针类型，
	// 此时GO会自动解析值类型为指针Exchange
	jerryll.Exchange()

	fmt.Println(jerryll)


	//结构内部为指针
	argPoint := "point-struct"
	stringPoint:=StringPoint{&argPoint}

	fmt.Println("raw stringPoint\t",stringPoint)

	argPoint ="changes"
	fmt.Println("chanage\tstringPoint\t",stringPoint)

	fmt.Println("chanage argPoint\t",argPoint)

	//raw stringPoint	 "point-struct"
	//chanage	stringPoint	 "changes"
	//chanage argPoint	 changes

	//结构内部为值
	argsValue := "value-struct"
	stringValue:=StringValue{argsValue}

	fmt.Println("raw\tstringValue\t",stringValue)

	argsValue ="changes"
	fmt.Println("change\tstringValue\t",stringValue)

	fmt.Println("change\targsValue",argsValue)

	//raw	stringValue	 "value-struct"
	//change	stringValue	 "value-struct"
	//change	argsValue changes
}