package safemap_test

import (
    "fmt"
    "safemap"
    "sync"
    "testing"
)

//
//WaitGroup的用途：它能够一直等到所有的goroutine执行完成，
// 并且阻塞主线程的执行，直到所有的goroutine执行完成。
//WaitGroup总共有三个方法：Add(delta int),Done(),Wait()。
//Add:添加或者减少等待goroutine的数量
//Done:相当于Add(-1)
//Wait:执行阻塞，直到所有的WaitGroup数量变成0

func TestSafeMap(t *testing.T) {
    store := safemap.New()
    fmt.Printf("Initially has %d items\n", store.Len())
    deleted := []int{0, 2, 3, 5, 7, 20, 399, 25, 30, 1000, 91, 97, 98, 99}
    //
    var waiter sync.WaitGroup
    waiter.Add(1)
    go func() { // Concurrent Inserter
        for i := 0; i < 100; i++ {
            store.Insert(fmt.Sprintf("0x%04X", i), i)
            if i > 0 && i%15 == 0 {
                fmt.Printf("Inserted %d items\n", store.Len())
            }
        }
        fmt.Printf("Inserted %d items\n", store.Len())
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Deleter
        for _, i := range deleted {
            key := fmt.Sprintf("0x%04X", i)
            before := store.Len()
            store.Delete(key)
            fmt.Printf("Deleted m[%s] (%d) before=%d after=%d\n",
                key, i, before, store.Len())
        }
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Finder
        for _, i := range deleted {
            for _, j := range []int{i, i + 1} {
                key := fmt.Sprintf("0x%04X", j)
                value, found := store.Find(key)
                if found {
                    fmt.Printf("Found m[%s] == %d\n", key, value)
                } else {
                    fmt.Printf("Not found m[%s] (%d)\n", key, j)
                }
            }
        }
        waiter.Done()
    }()
    waiter.Wait()

    //
    updater := func(value interface{}, found bool) interface{} {
        if found {
            return value.(int) * 1000
        }
        return 1
    }
    for _, i := range []int{5, 10, 15, 20, 25, 30, 35} {
        key := fmt.Sprintf("0x%04X", i)
        if value, found := store.Find(key); found {
            fmt.Printf("Original m[%s] == %d\t", key, value)
                store.Update(key, updater)
            if value, found := store.Find(key); found {
                fmt.Printf("Updated m[%s] == %5d\n", key, value)
            }
        }
    }

    fmt.Printf("Finished with %d items\n", store.Len())
    // not needed here but useful if you want to free up the goroutine
    data := store.Close()
    fmt.Println("Closed")
    fmt.Printf("len == %d\n", len(data))
    //for k, v := range data { fmt.Printf("%s = %v\n", k, v) }
}
