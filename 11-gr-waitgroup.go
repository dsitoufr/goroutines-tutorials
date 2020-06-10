package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func Sending(str string) <- chan string {
    ch := make(chan string)

       go func() {
            for i := 0; i < 10; i++ {
                ch <- fmt.Sprintf("%s, %d", str, i)
            }
       }()

    return ch
}

func main() {

    //sending to channel
    ch := Sending("joe")

    //receive from chan
    wg.Add(10)
    go func() {
        for {
            fmt.Println(<-ch)
            wg.Done()
        }
    }()

  //wait gr to terminate
wg.Wait()
fmt.Println("End of program..!")
}
