//从极坐标到笛卡儿坐标的转换
package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "runtime"
)

const result = "Polar radius=%.02f θ=%.02f° → Cartesian x=%.02f y=%.02f\n"

var prompt = "Enter a radius and an angle (in degrees), e.g., 12.5 90, " +
    "or %s to quit."

/**
*结构声明
 */
type polar struct {
    radius float64
    θ      float64
}

type cartesian struct {
    x   float64
    y   float64
}

/**
*init在main之前调用,但是不能显示调用
 */
func init() {
    if runtime.GOOS == "windows" {
        prompt = fmt.Sprintf(prompt, "Ctrl+Z, Enter")
    } else { // Unix-like
        prompt = fmt.Sprintf(prompt, "Ctrl+D")
    }
}

func main() {
    //通道创建声明与初始化,格式为chan Type
    //创建用来传输结构体类型为polar的通道,并赋于变量question
    questions := make(chan polar)
    //方法退出去关闭
    defer close(questions)
    //计算处理
    answers := createSolver(questions)
    defer close(answers)
    //计算转换处理
    interact(questions, answers)
}

/**
声明一个通道,入参为一个polar类型的通道,返回cartesian类型的通道
//创建一个字符串通道,通道的大小为10,通道类似于一个先进先出的队列
message:=make(chan string,10)
message<-"Leader"  二元操作时,左侧必须是通道,表示为通道赋值,通道满时则阻塞
message1:=<-message,一元操作符时,从通道中获取数据
 */
func createSolver(questions chan polar) chan cartesian {
    //创建用来传输结构体类型cartesian的通道
    answers := make(chan cartesian)
    //协程调用
    go func() {
        for {
            //一元操作符,从通道questions中获取数据
            polarCoord := <-questions
            θ := polarCoord.θ * math.Pi / 180.0 // degrees to radians
            x := polarCoord.radius * math.Cos(θ)
            y := polarCoord.radius * math.Sin(θ)
            //二元操作符,把数据赋值于通道,即把生成结构体值,并把它赋予结构体cartesian的通道

            answers <- cartesian{x, y}
        }
    }()
    return answers
}

func interact(questions chan polar, answers chan cartesian) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println(prompt)
    for {
        fmt.Printf("Radius and angle: ")
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        var radius, θ float64
        if _, err := fmt.Sscanf(line, "%f %f", &radius, &θ); err != nil {
            fmt.Fprintln(os.Stderr, "invalid input")
            continue
        }
        //赋值给结构体的通道
        questions <- polar{radius, θ}
        //从通道中取值
        coord := <-answers
        fmt.Printf(result, radius, θ, coord.x, coord.y)
    }
    fmt.Println()
}
