package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

/**
*命令行处理方法
 */
func handleCommandLine() (algorithm int, minSize, maxSize int64,
    suffixes, files []string) {
    flag.IntVar(&algorithm, "algorithm", 1, "1 or 2")
    flag.Int64Var(&minSize, "min", -1,
        "minimum file size (-1 means no minimum)")
    flag.Int64Var(&maxSize, "max", -1,
        "maximum file size (-1 means no maximum)")
    var suffixesOpt *string = flag.String("suffixes", "",
        "comma-separated list of file suffixes")
    flag.Parse()
    if algorithm != 1 && algorithm != 2 {
        algorithm = 1
    }
    if minSize > maxSize && maxSize != -1 {
        log.Fatalln("minimum size must be < maximum size")
    }
    suffixes = []string{}
    if *suffixesOpt != "" {
        suffixes = strings.Split(*suffixesOpt, ",")
    }
    files = flag.Args()
    return algorithm, minSize, maxSize, suffixes, files
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    log.SetFlags(0)
    //算法/文件最小值,最大值/文件后缀/文件列表
    algorithm,
        minSize, maxSize, suffixes, files := handleCommandLine()
    //通道
    if algorithm == 1 {
        sink(filterSize(minSize, maxSize, filterSuffixes(suffixes, source(files))))
    } else {
        //文件名切片,根据文件名切片,返回一个文件名的通道
        channel1 := source(files)
        //过滤符合后缀之后的文件名称通道
        channel2 := filterSuffixes(suffixes, channel1)
        //符合文件大小的通道
        channel3 := filterSize(minSize, maxSize, channel2)
        sink(channel3)
    }
}

//输出通道,创建1000个通道存储文件名，所以没有到达1000个时不阻塞
//当发送量占满整个缓冲区时阻塞，直到有空通道
//返回只读通道
//
func source(files []string) <-chan string {
    //带缓冲区通道
    out := make(chan string, 1000)
    //go语句构建匿名函数，构建之后立即执行。
    //因此调用source则协程开始执行
    //定义协程
    go func() {
        //轮询把文件名发送到通道out上
        for _, filename := range files {
            //通道接收，在没有到达通道容量时，不阻塞
            out <- filename
        }
        //处理完毕，由于是发送通道，需要显示关闭
        close(out)
    }()
    return out
}

// make the buffer the same size as for files to maximize throughput
//文件后缀的通道,通道大小和输入的通道一致
func filterSuffixes(suffixes []string, in <-chan string) <-chan string {
    //和入参in通道相同的缓冲区，增加吞吐量
    out := make(chan string, cap(in))
    //定义协程
    go func() {
        //轮询遍历in通道，获取文件名
        for filename := range in {//从in通道处获取数据
            //后缀长度为0，即没有指定后缀则全部接收文件
            if len(suffixes) == 0 {//文件长度为0时
                out <- filename
                continue
            }
            //如果指定了文件后缀
            ext := strings.ToLower(filepath.Ext(filename))
            //文件名过了通道,满足则终端
            for _, suffix := range suffixes {
                if ext == suffix {//当后缀相同时接收文件后跳出循环
                    out <- filename
                    break
                }
            }
        }
        //处理完毕，由于是发送通道，需要显示关闭
        close(out)
    }()
    //根据过滤器滤掉的文件名称通道
    return out
}

// make the buffer the same size as for files to maximize throughput
//最大值/最小值/文件名称,过滤文件,返回满足的通道
func filterSize(minimum, maximum int64, in <-chan string) <-chan string {
    out := make(chan string, cap(in))
    go func() {
        for filename := range in {//从通道中获取数据
            if minimum == -1 && maximum == -1 {
                out <- filename // don't do a stat call it not needed
                continue
            }
            finfo, err := os.Stat(filename)
            if err != nil {
                continue // ignore files we can't process
            }
            size := finfo.Size()
            if (minimum == -1 || minimum > -1 && minimum <= size) &&
                (maximum == -1 || maximum > -1 && maximum >= size) {
                out <- filename
            }
        }
        close(out)
    }()
    return out
}
//source、filterSuffixes、filterSize分别在自己的协程中执行
/**
*单通道,显示文件名
 */
func sink(in <-chan string) {
    for filename := range in {//变量只读通道，获取通道数据
        fmt.Println(filename)
    }
}

