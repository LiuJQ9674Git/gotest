//每个goroutine有自己的Map,而后合并此Map,即通过内存换取性能

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
)

var workers = runtime.NumCPU()

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    lines := make(chan string, workers*4)
    //以空间换取时间，每个协程有自己的map用于存储计算结果
    results := make(chan map[string]int, workers)
    //读取数据goroutine
    go readLines(os.Args[1], lines)
    getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
    //处理数据,处理结果通道
    for i := 0; i < workers; i++ {
        //并发处理的线程个数，计算之后把结果发生到results
        go processLines(results, getRx, lines)
    }
    totalForPage := make(map[string]int)
    //阻塞等待,合并数据,内存合并
    merge(results, totalForPage)
    //阻塞等待
    showResults(totalForPage)
}
func readLines(filename string, lines chan<- string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal("failed to open the file:", err)
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if line != "" {
            lines <- line
        }
        if err != nil {
            if err != io.EOF {
                log.Println("failed to finish reading the file:", err)
            }
            break
        }
    }
    close(lines)
}

func processLines(results chan<- map[string]int, getRx *regexp.Regexp,
    lines <-chan string) {
    //同并发的两个版本不同之处，每个协程创建本地的map来保存数据
    countForPage := make(map[string]int)
    for line := range lines {
        if matches := getRx.FindStringSubmatch(line); matches != nil {
            countForPage[matches[1]]++
        }
    }
    //处理之后把本地的处理结果(页处理后的map)写通道
    results <- countForPage
}
//从结果通道接收数据并合并后更新map
func merge(results <-chan map[string]int, totalForPage map[string]int) {
    for i := 0; i < workers; i++ {
        countForPage := <-results//取出数据
        for page, count := range countForPage {
            totalForPage[page] += count
        }
    }
}

func showResults(totalForPage map[string]int) {
    for page, count := range totalForPage {
        fmt.Printf("%8d %s\n", count, page)
    }
}
