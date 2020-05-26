package main

//
// How to use channel directions here 
// we test receive only channel
//

import (
    "fmt"
    "time"
)

func Sending(msg string) <- chan string {
   c := make(chan string)

   go func() {
      for i := 0; i < 5; i++ {
          c <- fmt.Sprintf("%s, %d", msg, i)
          time.Sleep(time.Millisecond)
       }
   }()

   return c

}

func Merge(input1, input2 <-chan string) <-chan string {
    
    c := make(chan string)

    go func() { 
        for { 
              c <- <-input1
        } 
    }()  

    go func() { 
        for { 
             c <- <-input2 
        } 
    }()

    return c
}

func main() {
    c := Merge(Sending("Joe"), Sending("Ann"))

   // for i := 0; i < 10; i++ {
   //     fmt.Println(<-c)
   // }
   // or
   // you can replace loop by go func() plus wait in main

    go func() {
        for {
            fmt.Println(<-c)
        }
    }()

    time.Sleep(10 * time.Millisecond)
    fmt.Println("Leaving mean program..!")
}
