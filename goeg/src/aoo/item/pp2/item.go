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
	*Item //嵌入
	catalogId int //聚合
}

type LuxuryItem struct {
	Item //嵌入
	makeup int //聚合
	dpMutex  *sync.RWMutex //需要定义传递
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
	item.dpMutex.RLock()
	defer item.dpMutex.RUnlock()
	//使用嵌入的方法
	return item.Item.Cost()*float64(item.makeup)
}

func main(){
	specialItem:=SpecialItem{&Item{"Green",20,5},1}
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,specialItem.catalogId)
	fmt.Println(specialItem.Cost())
	specialItem.SetPrice(40)
	fmt.Println(specialItem.id,specialItem.price,specialItem.quantity,specialItem.catalogId)
	fmt.Println(specialItem.Cost())
	//RWMutex指针类型需要声明
	dpMutex:=new(sync.RWMutex);
	//覆盖嵌入方法
	luxucyItem:=LuxuryItem{Item: Item{"Green",20,5}, makeup: 1,
		dpMutex: dpMutex}

	//luxucyItem:=LuxuryItem{Item: Item{"Green",20,5}, makeup: 1}
	fmt.Println(luxucyItem.Item.id,luxucyItem.price,luxucyItem.quantity,luxucyItem.makeup)
	fmt.Println(luxucyItem.Cost())


}
