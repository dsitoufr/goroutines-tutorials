package main 

import (
    "fmt"
)

func daisy(in, out chan int) {
    out <-  1+ <- in
}

const N=1000

func main() {
     
    //To avoid ovewritting in channel
    //use use another channel c 
     c := make(chan int)
     in := c

    for i := 0; i < N; i++ {
        out := make(chan int)
        go daisy(in,out)
     
        //copy value
        in = out
    }

    //To avoid ovewritting in channel
    //use use another channel c 
    // in < - 1 
    
     // c <- 1 
    // or
    //go func() {
    //    c <- 1
    //}() 
    // or

    go func(ch chan int) {
           ch <- 1
    }(c)

     //last value is in channel
     //fmt.Println(<- out)
     fmt.Println(<- in ) 

     //last value is in channel
     //fmt.Println(<- out)
     fmt.Println(<- in ) 

}
