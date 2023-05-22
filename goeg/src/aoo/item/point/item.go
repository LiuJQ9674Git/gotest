package main

import (
	"fmt"
	"sync"
)
//自定义结构，结构内部的值是否为指针，GO可以根据可寻址来自动识别
//是否需要修改和接收者是指针类型还是值类型相关。
//在GO中，通道、切片、映射、函数和接口使用make定义，返回一个指向特定类型值的引用。
//定义语法为 new(Type)等价于&Type{}
//下面结构SpecialItem中的Item为嵌入，只要是可寻址，GO可以解析地址，
// 这样Item是否为指针，均可以获得正确的结果
//type SpecialItem struct {
//	*Item //嵌入
//	catalogId int //聚合
//}
//或者
//type SpecialItem struct {
//	Item //嵌入
//	catalogId int //聚合
//}
type Item struct {
	id string //集合
	price float64 //集合
	quantity int //集合
}
type SpecialItem struct {
	*Item //嵌入
	catalogId int //聚合
}

type LuxuryItem struct {
	Item //嵌入
	makeup int //聚合
	dpMutex  *sync.RWMutex
}

//构造函数，返回值或者指针没有发现区别
func New(id string,price float64,quantity int)  *Item{
	return &Item{id:id,price:price,quantity:quantity}
}
//此方法的接收者可以为值或者为指针，由于不改变其变量的值，
// 只是只读，计算结果不影响
func (item Item)Cost()float64{
	return item.price*float64(item.quantity)
}

//此方法的接收者必须为指针，因为需要改变price变量的值
func (item *Item)SetPrice(price float64){
	 item.price=price;
}

//返回为指针和值，测试没有发现区别
func (item *Item)copy() *Item {
	return &Item{item.id,item.price,item.quantity}
	
}
//此方法的接收者可以为值或者为指针，由于不改变其变量的值，
// 只是只读，计算结果不影响
//覆盖嵌入方法
func (item LuxuryItem)Cost()float64{
	//使用嵌入的方法
	item.dpMutex.RLock()
	defer item.dpMutex.RUnlock()
	return item.Item.Cost()*float64(item.makeup)
}

func main(){
	//修改为值类型，结果相同
	specialItem:=SpecialItem{&Item{"Green",20,5},1}
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,specialItem.catalogId)
	fmt.Println(specialItem.Cost())
	fmt.Println("set price")
	specialItem.SetPrice(40)
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,specialItem.catalogId)
	fmt.Println(specialItem.Cost())
	dpMutex:=new(sync.RWMutex);
	//覆盖嵌入方法
	luxucyItem:=LuxuryItem{Item: Item{"Green",20,5}, makeup: 1,
	dpMutex: dpMutex}
	fmt.Println(luxucyItem.id,luxucyItem.price,luxucyItem.quantity,luxucyItem.makeup)
	fmt.Println(luxucyItem.Cost())
	fmt.Println("set price")
	luxucyItem.SetPrice(50)
	fmt.Println(luxucyItem.Cost())
	//构造方法
	item:=New("Red",10,30)
	fmt.Println(item.Cost())

	itemCopy:=item.copy();
	fmt.Println(itemCopy.Cost())
	itemCopy.SetPrice(90)
	fmt.Println(itemCopy.Cost())
	fmt.Println(item.Cost())
	//运行结果
}
