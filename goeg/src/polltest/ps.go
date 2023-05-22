package main

import "fmt"

func main() {
	fmt.Println("FD")
	fd, _:= New("test",true)
	fd.Init("test",false)
	println(fd.String())
	fmt.Println("fd.pd")
	fd.Close()
	println(fd.String())
	println(fd.pd.String())


	println(fd.pd.runtimeCtx)
}

type FD struct {

	pd pollDesc

	isFile bool

	name string
}

func (fd *FD) Init(name string, pollable bool) error {

	fd.isFile = pollable

	err := fd.pd.init(fd)

	return err
}

func (fd *FD) Close() error {
	fd.pd.evict()

	err := fd.evict()

	return err
}

func (fd *FD) evict() error {
	fd.pd.evict()
	return nil
}

func New(name string, isFile bool) (*FD, error){
	ret := &FD{
		name: name,
		isFile: isFile,
	}
	return ret, nil
}

func (fd FD) String() string{
	return fmt.Sprintf("%s %s %d",fd.name,fd.isFile,fd.pd)
	//return fmt.Sprintf(fd.name)
}

type pollDesc struct {
	runtimeCtx int16
}

func (pd *pollDesc) init(fd *FD) error {
	var ctx int16
	ctx  =1<<3
	pd.runtimeCtx = ctx
	return nil
}

func (pd *pollDesc) close() {
	pd.runtimeCtx = 1<<0
}

func (pd *pollDesc) evict() {
	pd.runtimeCtx=1<<2
}

func (pd pollDesc) String() string{
	return fmt.Sprintf("%s",pd.runtimeCtx)
	//return fmt.Sprintf(fd.name)
}