
package main

import (
"fmt"
"sync"
)

type Item struct {
	id string //集合
	price float64 //集合
	quantity int //集合
}
type SpecialItem struct {
	Item //嵌入
	catalogId int //聚合
}
type LuxuryItem struct {
	Item //嵌入
	makeup int //聚合
	dpMutex  sync.RWMutex //值类型无需声明
}

//此方法的接收者可以为值或者为指针，由于不改变其变量的值，
// 只是只读，计算结果不影响
func (item *Item)Cost()float64{
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
	item.dpMutex.RLock()
	defer 	item.dpMutex.RUnlock()
	result :=item.Item.Cost()*float64(item.makeup)

	return result
}

//构造函数，返回值或者指针没有发现区别
func New(id string,price float64,quantity int)  *Item{
	return &Item{id:id,price:price,quantity:quantity}
}

func main(){

	//修改为指针之后，同值类型
	specialItem:=&SpecialItem{Item{"Green",20,5},1}
	fmt.Println("修改之前")
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,
		specialItem.catalogId)//修改前的值
	fmt.Println(specialItem.Cost())
	specialItem.SetPrice(40)
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,
		specialItem.catalogId)//修改后的值
	fmt.Println(specialItem.Cost())
	fmt.Println("修改之后")
	//dpMutex:=new(sync.RWMutex);
	//覆盖嵌入方法
	luxucyItem:=LuxuryItem{Item: Item{"Green",20,5}, makeup: 1}
	fmt.Println(luxucyItem.Item.id,luxucyItem.price,luxucyItem.quantity,luxucyItem.makeup)
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