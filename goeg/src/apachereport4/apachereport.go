//

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
    "safeslice"
)

var workers = runtime.NumCPU()

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    lines := make(chan string)
    done := make(chan struct{}, workers)
    pageList := safeslice.New()
    go readLines(os.Args[1], lines)
    getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
    for i := 0; i < workers; i++ {
        //
        go processLines(done, getRx, pageList, lines)
    }
    waitUntil(done)
    showResults(pageList)
}

func readLines(filename string, lines chan<- string) {
    var file *os.File
    var err error
    if file, err = os.Open(filename); err != nil {
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

func processLines(done chan<- struct{}, getRx *regexp.Regexp,
    pageList safeslice.SafeSlice, lines <-chan string) {
    for line := range lines {
        if matches := getRx.FindStringSubmatch(line); matches != nil {
            pageList.Append(matches[1])
        }
    }
    done <- struct{}{}
}

func waitUntil(done <-chan struct{}) {
    for i := 0; i < workers; i++ {
        <-done
    }
}

func showResults(pageList safeslice.SafeSlice) {
    list := pageList.Close()
    counts := make(map[string]int)
    for _, page := range list { // uniquify
        counts[page.(string)] += 1
    }
    for page, count := range counts {
        fmt.Printf("%8d %s\n", count, page)
    }
}
