//////////////////////////////////////
// Goroutine sync with:
// WaitGroup
// Channel
// Channel message confirmation
//////////////////////////////////////////



package main 

import (
    "fmt"
    "sync"
    "runtime"
)

func Sender( str string, strChan chan string, quitChan chan bool) {

   i := 0
   for {
       select {
             case  strChan <- fmt.Sprintf("%s,%d",str, i):
                    i++
             case  <- quitChan:
                  msgChan <- "Bye Bye...!"
                  break
       }
   }
}

var (
    wg sync.WaitGroup
    msgChan = make(chan string)
)

func main() {

    strChan := make(chan string)
    quitChan := make(chan bool)
     
     //sender
     go Sender("joe", strChan, quitChan)
    
     //receiver: loop channel
     wg.Add(20) 
     go func() {
          for {
              select {
                   case s := <- strChan:
                        fmt.Println(s)
                        wg.Done()
              }
          }
     }()

     //wait GR
      wg.Wait()

     //until...
     quitChan <- true
     fmt.Println(<- msgChan, runtime.NumGoroutine())

}