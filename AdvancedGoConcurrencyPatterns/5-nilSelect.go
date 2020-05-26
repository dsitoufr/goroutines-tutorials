package main

import (
    "fmt"
    "math/rand"
    "time"
)

func init() {
  rand.Seed(time.Now().UnixNano())
}

func main() {

    a := make(chan string)
    b := make(chan string)

    go func(){
        a <- "a"
    }()

    go func() {
        b <- "b"
    }()

    if rand.Intn(2) == 0 {
        a = nil
        fmt.Println("a nil")
    } else {
        b = nil 
        fmt.Println("b nil")
    }

    select {
        case <- a:
              fmt.Println("got a")
        case <- b:
              fmt.Println("got b")
        //default:  fmt.Println("nothing")  //if add default it will always hit default
              
    }

}