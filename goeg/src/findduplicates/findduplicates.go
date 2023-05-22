//filepath标准库来遍历路径
//程序根据工作量的多少来决定使用多少个goroutine
//对于大文件会有一个goroutine被单独创建以用于计算文件的SHA-1值,
//而小文件则是直接在当前goroutine计算
//也就是说我们不知道到底有多少个goroutine,但是需要设置一个上限

package main

import (
    "crypto/sha1"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "sort"
    "sync"
)

const maxSizeOfSmallFile = 1024 * 32
//
const maxGoroutines = 100

//路径信息
type pathsInfo struct {
    size  int64
    paths []string
}

//文件信息
type fileInfo struct {
    sha1 []byte
    size int64
    path string
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <path>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    //文件通道,通道的大小为200
    infoChan := make(chan fileInfo, maxGoroutines*2)
    //子goroutine出来找的过程
    go findDuplicates(infoChan, os.Args[1])
    //主goroutine合并结果为路径
    pathData := mergeResults(infoChan)
    //主goroutine输出结果
    outputResults(pathData)
}

func findDuplicates(infoChan chan fileInfo, dirname string) {
    //同步锁
    waiter := &sync.WaitGroup{}
    //遍历路径
    filepath.Walk(dirname, makeWalkFunc(infoChan, waiter))
    //阻塞直到操作完成
    waiter.Wait() // Blocks until all the work is done
    close(infoChan)
}

//遍历路径的具体函数
func makeWalkFunc(infoChan chan fileInfo,
    waiter *sync.WaitGroup) func(string, os.FileInfo, error) error {
    return func(path string, info os.FileInfo, err error) error {
        if err == nil && info.Size() > 0 &&
            (info.Mode()&os.ModeType == 0) {
            //文件小于边界尺寸直接计算,或者运行时的goroutine数量大于边界则直接计算
            if info.Size() < maxSizeOfSmallFile ||
                runtime.NumGoroutine() > maxGoroutines {
                processFile(path, info, infoChan, nil)
            } else {//否则启动新的goroutine计算
                waiter.Add(1)
                go processFile(path, info, infoChan,
                    func() {
                        waiter.Done() //互斥锁表示做完
                    })
            }
        }
        return nil // We ignore all errors
    }
}

func processFile(filename string, info os.FileInfo,
    infoChan chan fileInfo, done func()) {
    if done != nil {
        defer done()//最后执行done方法
    }
    file, err := os.Open(filename)
    if err != nil {
        log.Println("error:", err)
        return
    }
    defer file.Close()//方法退出时关闭文件
    //sha1实例
    hash := sha1.New()
    if size, err := io.Copy(hash, file);
        size != info.Size() || err != nil {
        //
        if err != nil {
            log.Println("error:", err)
        } else {
            log.Println("error: failed to read the whole file:", filename)
        }
        return
    }
    //文件信息写入文件通道
    infoChan <- fileInfo{hash.Sum(nil), info.Size(), filename}
}

func mergeResults(infoChan <-chan fileInfo) map[string]*pathsInfo {
    pathData := make(map[string]*pathsInfo)
    format := fmt.Sprintf("%%016X:%%%dX", sha1.Size*2) // == "%016X:%40X"
    for info := range infoChan {//读取文件信息通道数据
        key := fmt.Sprintf(format, info.size, info.sha1)
        value, found := pathData[key]
        if !found {
            //新建paths信息
            value = &pathsInfo{size: info.size}
            pathData[key] = value
        }
        value.paths = append(value.paths, info.path)
    }
    return pathData
}

func outputResults(pathData map[string]*pathsInfo) {
    keys := make([]string, 0, len(pathData))
    for key := range pathData {
        keys = append(keys, key)
    }
    sort.Strings(keys)
    for _, key := range keys {
        value := pathData[key]
        if len(value.paths) > 1 {
            fmt.Printf("%d duplicate files (%s bytes):\n",
                len(value.paths), commas(value.size))
            sort.Strings(value.paths)
            for _, name := range value.paths {
                fmt.Printf("\t%s\n", name)
            }
        }
    }
}

// commas() returns a string representing the whole number with comma
// grouping.
func commas(x int64) string {
    value := fmt.Sprint(x)
    for i := len(value) - 3; i > 0; i -= 3 {
        value = value[:i] + "," + value[i:]
    }
    return value
}
