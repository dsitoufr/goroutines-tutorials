package main

import (
    "fmt"
//"time"
)

func main() {

    bidirectChan := make(chan string)

   // Receiver
    go func() {
        fmt.Println("waiting channel..")
        msg := <- bidirectChan 
        fmt.Println("got channel: ",  msg)
    }()

     //sender
    fmt.Println("sending channel...")
    bidirectChan <- "Hello"
    fmt.Println("End program.")
    

}