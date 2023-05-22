package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "runtime"
)

var workers = runtime.NumCPU()

type Result struct {
    filename string
    lino     int
    line     string
}

type Job struct {
    filename string
    //写入通道，即单向只允许发送数据通道
    results  chan<- Result
}

////////////////////////////
func main() {
    //使用的最大核数
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <regexp> <files>\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    //regexp指针类型Regexp的变量,传给grep函数被所有工作的的goroutine共享.
    //注意,一般的我们必须假设任何共享指针指向的值都是线程不安全的.
    //对于Regexp,go语言中,这个指针指向的值是线程安全
    //regexp.Compile为指针类型,*Regexp线程安全
    if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
        log.Fatalf("invalid regexp: %s\n", err)
    } else {
        //
        grep(lineRx, commandLineFiles(os.Args[2:]))
    }
}

//获取文件切片
func commandLineFiles(files []string) []string {
    if runtime.GOOS == "windows" {
        args := make([]string, 0, len(files))
        for _, name := range files {
            if matches, err := filepath.Glob(name); err != nil {
                args = append(args, name) // Invalid pattern
            } else if matches != nil { // At least one match
                args = append(args, matches...)
            }
        }
        return args
    }
    return files
}

//入参表达式和文件切片
func grep(lineRx *regexp.Regexp, filenames []string) {
    //三个带有缓冲区通道，其中jobs和done缓冲为机器的CPU梳理
    jobs := make(chan Job, workers)//Job双通道,大小和核数一致,即任务大小为核数
    //Result双通道,大小为1000和文件数量的最小值
    results := make(chan Result, minimum(1000, len(filenames)))
    done := make(chan struct{}, workers) //处理过程控制通道,done通道和任务job的通道数量一样
    //在自己的(子协程中)goroutine执行,它和主线程并行执行
    go addJobs(jobs, filenames, results) // Executes in its own goroutine
    //启动工作协程梳理
    for i := 0; i < workers; i++ {//以每个核作为一次处理,处理Job
        //在子协程中处理执行任务,
        go doJobs(done, lineRx, jobs) // Each executes in its own goroutine
    }
    go awaitCompletion(done, results) // Executes in its own goroutine

    processResults(results)//阻塞直到任务处理完成,消费数据Blocks until the work is done
}

//协程：jobs和results为发送通道，只允许发数据
//for，select关闭通道，或者使用->用来检查是否可以接收时才显示关闭通道
//应该由发送端的线程关闭通道，而不是由接收端线程完成
//即使没有检查通道关闭，对性能的影响也是轻微的

//jobs通道在for/select结构中，所以由发送端关闭
//done通道没有再for/select中，所以无需显示关闭
func addJobs(jobs chan<- Job, filenames []string, results chan<- Result) {
    for _, filename := range filenames {//filenames为jobList
        //如果任务的大小满,则阻塞,等待接收,否则把Job数据放入job通道
        //阻塞等待接收，创建类型Job
        jobs <- Job{filename, results}
    }
    //关闭通道，接收者可以结束
    close(jobs)
}

//协程：jobs为接收通道，只允许通道中获取数据
func doJobs(done chan<- struct{}, lineRx *regexp.Regexp, jobs <-chan Job) {
    for job := range jobs {//无任务阻塞,等待发送任务,
                            // 有则接收任务,
        job.Do(lineRx)//处理任务，Do中在results通道中写入数据
    }
    //任务做完标识，阻塞，等待接收
    done <- struct{}{}
}

////////////////////////
func (job Job) Do(lineRx *regexp.Regexp) {
    file, err := os.Open(job.filename)
    if err != nil {
        log.Printf("error: %s\n", err)
        return
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for lino := 1; ; lino++ {
        line, err := reader.ReadBytes('\n')
        line = bytes.TrimRight(line, "\n\r")
        if lineRx.Match(line) {
            //向通道中写入数据，发送通道
            job.results <- Result{job.filename,
                lino, string(line)}
        }
        if err != nil {
            if err != io.EOF {
                log.Printf("error:%d: %s\n", lino, err)
            }
            break
        }
    }
}

//协程：done消费之后，关闭结果通道results
func awaitCompletion(done <-chan struct{}, results chan Result) {
    for i := 0; i < workers; i++ {
        <-done//等待任务做完,如果任务处理,
              // 则从任务处理的Done通道中取出数据,表示已经处理完
    }
    close(results)
}
//results为接收通道，只允许通道中获取数据
func processResults(results <-chan Result) {
    for result := range results {//如果Result通道中没有数据,则阻塞等待数据发送,
                                // 如果有则直接处理数据
        fmt.Printf("%s:%d:%s\n", result.filename,
            result.lino, result.line)
    }
}

func minimum(x int, ys ...int) int {
    for _, y := range ys {
        if y < x {
            x = y
        }
    }
    return x
}

