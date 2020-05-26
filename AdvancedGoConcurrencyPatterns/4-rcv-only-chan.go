package main

import (
    "fmt"
)

func main() {

    // use bd to init value
    bidirectChan := make(chan string)

    // sender
    go func() {
        bidirectChan <- "Hello world"
    }()
    
    // receiver
    receiveOnlyFunc(bidirectChan)

    fmt.Println("End program.")

}

func receiveOnlyFunc(  c <- chan string) {
  fmt.Println( <- c ) 
}