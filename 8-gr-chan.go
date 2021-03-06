package main

//
// How to control a GR with a boolean
//

import (
    "fmt"
    "time"
)

type Message struct {
    str string
    wait chan bool
}

func Watching(msg string) <- chan Message {

  waitChan := make(chan bool)
  ch := make(chan Message)

  go func() {
   for i:= 0; ; i++ {
      ch  <-Message { str: fmt.Sprintf("%s,%d",msg,i),  wait: waitChan }
      time.Sleep(time.Millisecond)
      <- waitChan
    }
  }()

  return ch
}

func allfan(a, b <- chan Message) <- chan Message {
    ch := make(chan Message)

    /*
    go func() {
        for {  //loop chan a
             ch <- <- a
        }
    }()

    go func() {  //loop chan b
       for {
           ch <- <- b
       }
    }()
    
    OR */

   go func() {
       for {
           select {
               case  s := <- a:  ch <- s
               case  s := <- b: ch <- s
           }
       }
   }()

    return ch
}

func main() {

    ch := allfan( Watching("joe"),  Watching("jack"))

    
   for i:= 0; i < 10; i++ {
        msg1 := <- ch
        fmt.Println(msg1.str)
        msg2 := <- ch
        fmt.Println(msg2.str)

        msg1.wait <- true
        msg1.wait <- true
    }

    fmt.Println("End program.")
}
