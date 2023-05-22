package main

import (
    "bufio"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

var britishAmerican = "british-american.txt"

func init() {
    dir, _ := filepath.Split(os.Args[0])
    britishAmerican = filepath.Join(dir, britishAmerican)
}

func main() {
    inFilename, outFilename, err := filenamesFromCommandLine()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    inFile, outFile := os.Stdin, os.Stdout
    if inFilename != "" {
        if inFile, err = os.Open(inFilename); err != nil {
            log.Fatal(err)
        }
        //任何属于defer语句所对应的语句都会被执行
        //但是此语句只会在使用defer的函数在返回时被调用
        //其执行的控制权会马上交给下一个语句
        defer inFile.Close()
    }
    if outFilename != "" {
        if outFile, err = os.Create(outFilename); err != nil {
            log.Fatal(err)
        }
        defer outFile.Close()
    }

    if err = americanise(inFile, outFile); err != nil {
        log.Fatal(err)
    }
}

/**
*具名返回值,作用域
 */
func filenamesFromCommandLine() (inFilename, outFilename string,
    err error) {
    if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
        err = fmt.Errorf("usage: %s [<]infile.txt [>]outfile.txt",
            filepath.Base(os.Args[0]))
        return "", "", err
    }
    if len(os.Args) > 1 {
        inFilename = os.Args[1]
        if len(os.Args) > 2 {
            outFilename = os.Args[2]
        }
    }
    if inFilename != "" && inFilename == outFilename {
        log.Fatal("won't overwrite the infile")
    }
    return inFilename, outFilename, nil
}

//具名返回值,不应当使用:=来赋值
//
func americanise(inFile io.Reader, outFile io.Writer) (err error) {
    reader := bufio.NewReader(inFile)
    writer := bufio.NewWriter(outFile)
    //err为具名返回,必须保证不使用快速声明变量:=,以避免影子变量
    //匿名的延迟函数,在函数返回时调用
    defer func() {
        if err == nil {
            err = writer.Flush()
        }
    }()

    //声明时赋值,避免影子变量
    //基于这种考虑,必须先声明变量,如replacer变量/line变量的什么
    //if :=会创建一个一个新的块,它会隐藏同名变量的值
    var replacer func(string) string
    if replacer, err = makeReplacerFunction(britishAmerican); err != nil {
        return err
    }
    //返回一个Regexp的指针
    wordRx := regexp.MustCompile("[A-Za-z]+")
    eof := false
    for !eof {
        var line string
        //读UTF-8的编码的字节流
        line, err = reader.ReadString('\n')
        if err == io.EOF {
            err = nil   // io.EOF isn't really an error,并不是一个真实的error
            eof = true  // this will end the loop at the next iteration
        } else if err != nil {
            return err  // finish immediately for real errors
        }
        //ReplaceAllStringFunc,接收一个字符串,和一个函数,如果匹配则执行函数,
        // 并将匹配的内容代替为replacer函数的返回内容
        line = wordRx.ReplaceAllStringFunc(line, replacer)
        //WriteString写UTF-8的字节流
        if _, err = writer.WriteString(line); err != nil {
            return err
        }
    }
    return nil
}

/**
*返回Replacer函数,返回的函数具有两个返回值,一个函数,一个是错误码
*待处理的文件和用来代替的字符串
 */
func makeReplacerFunction(file string) (func(string) string, error) {
    rawBytes, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    text := string(rawBytes)
    //切片/映射/数组使用make创建,分配类型大小的内存并初始化,指向特定值引用
    usForBritish := make(map[string]string)
    lines := strings.Split(text, "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) == 2 {
            usForBritish[fields[0]] = fields[1]
        }
    }

    return func(word string) string {
        if usWord, found := usForBritish[word]; found {
            return usWord
        }
        return word
    }, nil
}
