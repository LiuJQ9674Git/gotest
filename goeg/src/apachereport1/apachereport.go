//并发处理需要解决更新共享数据问题,在Go中有两种方案,一是使用互斥锁,二是使用通道达到串行化
//使用通道和安全映射(SafeMap),通过互斥量保护共享的map,还有通过通道局部的map来避免访问
//串行化从而达到吞吐量最大化,并使用通道对一个map进行更新

//功能:读取access.log,然后统计所有记录里访问html的次数
//实现:用一个goroutine读取日志中的一行,另外三个goroutine处理,每读到一个html被访问时,就更新到映射中.
//尽管多个goroutine同时处理这些记录,但是他们所有的更新使用同一个map

package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "runtime"
    "safemap"
)

var workers = runtime.NumCPU()

func main() {
    //使用机器的所有核数
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    //行的通道为string类型,大小为核数的4倍
    lines := make(chan string, workers*4)
    //处理情况的保存,大小为核数
    done := make(chan struct{}, workers)
    //创建safemap通道,safemap为非导出类型的值,实现了SafeMap,可以任意传递
    pageMap := safemap.New()
    //启动一个goroutine从日志文件中读取数据,即读取行goroutine
    go readLines(os.Args[1], lines)
    //启动处理日志中一行数据子goroutine
    processLines(done, pageMap, lines)
    //等待所有任务处理完成
    waitUntil(done)
    //显示结果
    showResults(pageMap)
}

func readLines(filename string, lines chan<- string) {//只允许发送到通道
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal("failed to open the file:", err)
    }
    defer file.Close()
    //lines通道设置小缓冲区,降低goroutine阻塞在lines通道上的可能
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if line != "" {//写入数据到通道
            lines <- line
        }
        if err != nil {
            if err != io.EOF {
                log.Println("failed to finish reading the file:", err)
            }
            break
        }
    }
    //关闭缓冲区
    close(lines)
}

func processLines(done chan<- struct{}, pageMap safemap.SafeMap,
    lines <-chan string) {
    getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
    //SafeMap中的值值被定义为接口类型,所以需要,所以需要在此方法中做类型的断言
    incrementer := func(value interface{}, found bool) interface{} {
        if found {
            return value.(int) + 1//类型断言
        }
        return 1
    }
    for i := 0; i < workers; i++ {
        //workers为具体处理的通道数量
        go func() {
            //每个行作为结果
            for line := range lines {//lines为只允许发送数据的通道,当满时,则阻塞,等待消费.
                if matches := getRx.FindStringSubmatch(line);
                    matches != nil {//更新数据
                    pageMap.Update(matches[1], incrementer)
                }
            }
            //更新完毕把处理步骤写入done通道
            done <- struct{}{}
        }()
    }
}

func waitUntil(done <-chan struct{}) {
    for i := 0; i < workers; i++ {
        <-done
    }
}

func showResults(pageMap safemap.SafeMap) {
    pages := pageMap.Close()
    for page, count := range pages {
        fmt.Printf("%8d %s\n", count, page)
    }
}
