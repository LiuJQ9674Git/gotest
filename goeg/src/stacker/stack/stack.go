package stack

import "errors"
//接口也是一种类型,空接口可以表示任意类型的值,在Go中仅仅有类型和值
//Stack可以包含任意类型混合的存储
//Go是强制类型,使用接口interface{}可以定义任何类型
type Stack []interface{}

//函数和方法的区别,方法需要定义所作用域的类型
//调用该方法的值,即接收器
//在指向的Stack指针上进行调用
//如果需要修改接收器,必须将接收器定义为一个指针
func (stack *Stack) Pop() (interface{}, error) {
    theStack := *stack
    if len(theStack) == 0 {
        return nil, errors.New("can't Pop() an empty stack")
    }
    x := theStack[len(theStack)-1]
    *stack = theStack[:len(theStack)-1]
    return x, nil
}

func (stack *Stack) Push(x interface{}) {
    *stack = append(*stack, x)
}

//接收器类型是按值传递的,任何对接收器的改变,仅仅改变其副本
//按值传递的,在不改变接收器副本的情况使用
func (stack Stack) Top() (interface{}, error) {
    if len(stack) == 0 {
        return nil, errors.New("can't Top() an empty stack")
    }
    return stack[len(stack)-1], nil
}


func (stack Stack) Cap() int {
    return cap(stack)
}

func (stack Stack) Len() int {
    return len(stack)
}

func (stack Stack) IsEmpty() bool {
    return len(stack) == 0
}
