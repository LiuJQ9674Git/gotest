package main

import "fmt"

type BenchmarkOptions struct {
	masters          *string //主地址，默认localhost:9333
	concurrency      *int //读协程并发数，默认16
	numberOfFiles    int //每个线程写文件数量
	fileSize         int //文件大小
}

//变量
var (
	b  BenchmarkOptions
)

func init()  {
	initDB()
}
func initDB() {
	master:="master"
	b.masters = &master
	concurrency:=16
	b.concurrency = &concurrency
	//numberOfFiles:=1024
	b.numberOfFiles=1024;

}

func printB() {
	fmt.Println(*b.concurrency,*b.masters,b.numberOfFiles,b.fileSize)

}

func printAssign(bb  *BenchmarkOptions) {
	fmt.Println(*bb.concurrency,*bb.masters,bb.numberOfFiles,bb.fileSize)

}
func main()  {
	printB()
	b.fileSize=2048
	master:="mastertwo"
	numberOfFiles:=1000
	b.masters = &master
	b.numberOfFiles=numberOfFiles
	printAssign(&b)
	//由于b.masters中为指针，所以变量值修改
	master="tresspoint"
	//由于b.masters中为值，所以变量值不修改
	numberOfFiles=2000
	printAssign(&b)
}