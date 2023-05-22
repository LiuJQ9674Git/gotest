package main

import (
	"fmt"
	"strconv"
)

type Exachanger interface {Exachange()}
type StringPair struct {firstName,lastName string}
type FloatPair struct {first,last float64}

func (stringPair *StringPair)Exachange()  {
	fmt.Println("StringPair Exachange...")
}
//String方法接收器为StringPair指针
func (stringPair *StringPair)String()string  {
	return stringPair.firstName+" test "+stringPair.lastName
}
func (floattPair FloatPair)Exachange()  {
	fmt.Println("FloatPair Exachange...")
}
func (floattPair FloatPair)String()string  {
	return  strconv.FormatFloat(floattPair.first,'f',-1,
		64)+","+strconv.FormatFloat(floattPair.last,'f',-1,
		64)
}
func handleExachanger(exachangers...Exachanger){
	//使用函数处理
	fmt.Println("处理Exachanger业务")
	for _,exachanger :=range exachangers{
		exachanger.Exachange()
	}
}

func main() {
	stringPair :=StringPair{"Bei","Jing"}
	stringPair.Exachange()
	//Println方法的形参为：
	//Println(a ...interface{}) (n int, err error)
	//因此不能输出test，而只能输出{}
	fmt.Println("StringPair is:",stringPair)
	floatPair :=FloatPair{110.00,222}
	floatPair.Exachange()
	//可以输出正确信息
	fmt.Println("FloatPair is:",floatPair)
	//函数处理
	handleExachanger(&stringPair,floatPair)
	testInferfaceAndStruct(stringPair,floatPair)
}

//形参为值
func testInferfaceAndStruct(stringPair StringPair,floatPair FloatPair)  {
	//类型断言
	var i interface{}=99
	if ii,ok :=i.(int) ; ok{//类型断言
		fmt.Println("断言 is:",ii)
	}
	var s interface{}=[]string{"first","last"}
	if ss,ok :=s.(string) ; ok{//类型断言
		fmt.Println("字符串断言值:",ss)
	}
	//泛型接口类型切片
	var ps interface{}=&StringPair{"first","last"}
	if p,found :=ps.(Exachanger) ; found{//类型断言
		fmt.Println("interface自定义的接口断言正确的结果:",p)//执行
	}else{
		fmt.Println("interface自定义的接口断言失败:",p)
	}
	//具体接口类型切片
	var psExachanger Exachanger=
		&StringPair{"first","last"}

	if p,found :=psExachanger.(Exachanger) ; found{//类型断言
		fmt.Println("Exachanger自定义的接口断言正确的结果:",p)//执行
	}else{
		fmt.Println("Exachanger自定义的接口断言失败:",p)
	}
	pStringPair :=&StringPair{"上海","中国"}
	var pp Exachanger=pStringPair
	if p,found :=pp.(Exachanger) ; found{
		fmt.Println("引用自定义的接口断言正确的结果:",p)//执行
	}else{
		fmt.Println("引用自定义的接口断言失败:",p)
	}
	//------------------------------------指针类型，值类型
	var pairs []Exachanger =[]Exachanger{&stringPair,floatPair}
	fmt.Println("Exachanger切片:",pairs)
	for i:=range pairs{
		fmt.Println("索引循环:",pairs[i])
	}
	for index, i:=range pairs{
		fmt.Println("内容循环:",index,i)
	}
	//String方法的接收器类型（值与指针），执行
	for _, ex :=range []Exachanger{pStringPair,floatPair}{
		fmt.Println("for中的自定义接口断言:",ex)
		ex.Exachange()
	}
}
