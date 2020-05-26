package main

import (
    "fmt"
    "time"
)

func watch(msg  string) <- chan string {

    c := make(chan string)
    go func() {
      for i := 0; i < 5; i++ {
          c <- fmt.Sprintf("%s, %d", msg, i)
          time.Sleep(10 * time.Millisecond)
       }
    }()

    return c
}

func main() {

    j := watch("joe")
    k := watch("ken")

    for i := 0; i < 5; i++ {
        fmt.Println( <- j)
        fmt.Println( <- k)
    }

    fmt.Println("End program...!")
}
