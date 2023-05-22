//线程安全的映射map,键为字符串,值为空接口
//如果值存的是指针或者引用,则必须保证指针指向的值是只读的,
//或者对于它们的访问是串行的
//
//<-将一个数据value写入至channel，
// 这会导致阻塞，直到有其他goroutine从这个channel中读取数据
// 将一个数据value写入至channel，
// 这会导致阻塞，直到有其他goroutine从这个channel中读取数据
//ch <- value

//单向channel变量的声明：
//chan为关键字
//var ch1 chan int  　　　　// 普通channel
//var ch2 chan <- int 　　 // 只用于写数据到int类型通道
//var ch3 <-chan int 　　 // 只用于从int类型通道读数据
//
//
//可以通过类型转换，将一个channel转换为单向的：
//
//ch4 := make(chan int)
//ch5 := <-chan int(ch4)   // 单向读
//ch6 := chan<- int(ch4)  //单向写

//var ch1 chan int  　　　　// 普通channel
//var ch2 chan <- int 　　 // 只用于写数据到int类型通道
//var ch3 <-chan int 　　 // 只用于从int类型通道读数据
package safemap

type commandAction int

const (
    remove commandAction = iota
    end
    find
    insert
    length
    update
)

//不可导出的结构
type commandData struct {
    action  commandAction//主要执行的动作action
    key     string
    value   interface{}
    result  chan<- interface{}//只写通道,只接收数据,不发送数据
    data    chan<- map[string]interface{}//只写通道
    updater UpdateFunc
}


//定义safeMap通道
type safeMap chan commandData

//不可导出的结构
type findResult struct {
    value interface{}
    found bool
}

//可导出的接口
type SafeMap interface {
    Insert(string, interface{})
    Delete(string)
    Find(string) (interface{}, bool)
    Len() int
    Update(string, UpdateFunc)
    Close() map[string]interface{}
}

type UpdateFunc func(interface{}, bool) interface{}

func New() SafeMap {
    //type safeMap chan commandData commandData通道声明
    sm := make(safeMap)
    //sm启动子goroutine运行,将safeMap同goroutine关联,sm在自己的goroutine中运行
    go sm.run()
    //返回sm,此sm在客户端使用Insert/Find/Update可安全的执行操作，即读取sm通道中的数据
    return sm
}

func (sm safeMap) run() {
    //保存安全map的key-value的map结构，在子协程中执行
    store := make(map[string]interface{})
    //遍历sm(safeMap)通道，获取定义的6种操作，如果通道是空则阻塞
    for command := range sm {//消费sm中的数据,
                            // commandData通道处理,sm如果为空则阻塞
        switch command.action {//消费数据
        case insert:
            //插入数据时,映射的store的key-value存储来自sm通道中的key-value
            store[command.key] = command.value
        case remove:
            //从保存的映射中删除
            delete(store, command.key)
        case find:
            //查找的结果
            value, found := store[command.key]
            //value与是否发现value作为数据写入只写通道
            command.result <- findResult{value, found}
        case length:
            //长度数据写入只写通道
            command.result <- len(store)
        case update:
            value, found := store[command.key]
            //updater函数调用safeMap方法（即这里的command）会导致死锁。
            //在command.updater函数不返回，则update分支不能正常结束时，
            //如果command.updater调用了一个safeMap方法（即这里的command）
            //它会一直阻塞到update分支完成，这样两个方法都完成不了导致死锁。
            store[command.key] = command.updater(value, found)
        case end:
            //数据写入只写通道
            close(sm)//关闭通道，run中的for退出
            command.data <- store//数据写入data通道
        }
    }
}

func (sm safeMap) Insert(key string, value interface{}) {
    //通道sm接收数据
    sm <- commandData{action: insert, key: key, value: value}
}

func (sm safeMap) Delete(key string) {
    //当客户端调用Delete时,发送数据commandData到sm通道中,
    //在New函数创建的sm所在的goroutine中消费sm
    sm <- commandData{action: remove, key: key}
}

func (sm safeMap) Find(key string) (value interface{}, found bool) {
    //定义接口类型的value通道
    reply := make(chan interface{})
    //通道接收数据
    sm <- commandData{action: find, key: key, result: reply}
    result := (<-reply).(findResult)//从通道reply中获取数据
    return result.value, result.found
}

func (sm safeMap) Len() int {
    reply := make(chan interface{})
    //通道接收数据
    sm <- commandData{action: length, result: reply}
    return (<-reply).(int)
}

// If the updater calls a safeMap method we will get deadlock!
func (sm safeMap) Update(key string, updater UpdateFunc) {
    //通道写入数据
    sm <- commandData{action: update, key: key, updater: updater}
}

// Close() may only be called once per safe map; all other methods can be
// called as often as desired from any number of goroutines
//关闭safeMap通道，这样就不会再有其它的更新
//关闭通道，run中的for退出，释放资源`
func (sm safeMap) Close() map[string]interface{} {
    //reply为映射类型的通道
    reply := make(chan map[string]interface{})
    //通道写入数据
    sm <- commandData{action: end, data: reply}
    return <-reply//将数据返给客户端
}

