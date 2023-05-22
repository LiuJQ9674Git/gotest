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
    results  chan<- Result
}

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
            job.results <- Result{job.filename, lino, string(line)}
        }
        if err != nil {
            if err != io.EOF {
                log.Printf("error:%d: %s\n", lino, err)
            }
            break
        }
    }
}

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

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <regexp> <files>\n",
            filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
        log.Fatalf("invalid regexp: %s\n", err)
    } else {
        grep(lineRx, commandLineFiles(os.Args[2:]))
    }
}

func grep(lineRx *regexp.Regexp, filenames []string) {
    jobs := make(chan Job, workers)
    results := make(chan Result, minimum(1000, len(filenames)))
    done := make(chan struct{}, workers)
    go addJobs(jobs, filenames, results)
    for i := 0; i < workers; i++ {
        go doJobs(done, lineRx, jobs)
    }
    //结果的处理合并
    waitAndProcessResults(done, results)
}

func addJobs(jobs chan<- Job, filenames []string, results chan<- Result) {
    for _, filename := range filenames {
        jobs <- Job{filename, results}
    }
    close(jobs)
}

func doJobs(done chan<- struct{}, lineRx *regexp.Regexp, jobs <-chan Job) {
    for job := range jobs {
        job.Do(lineRx)
    }
    done <- struct{}{}
}

func waitAndProcessResults(done <-chan struct{}, results <-chan Result) {
    for working := workers; working > 0; {//影子变量
        select { // Blocking，协程selete阻塞，没有default
        case result := <-results://从results消费
            fmt.Printf("%s:%d:%s\n", result.filename, result.lino,
                result.line)
        case <-done://从done消费
            working--
        }
    }
    //是否存在done已经消费，但是结果没有处理完毕的情况
DONE:
    for {
        select { // Nonblocking
        case result := <-results: //不阻塞，是因为有缺省的break
            fmt.Printf("%s:%d:%s\n", result.filename, result.lino,
                result.line)
        default:
            break DONE
        }
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

