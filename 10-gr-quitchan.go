////////////////////////////////////////////
// Goroutine Synhronization with
// WaitGroup
// Channel
//////////////////////////////////////////

package main

import(
    "fmt"
    "sync"
)

//sending
func Sending(str string, quit chan bool)  <- chan string {

    ch := make(chan string)
    i := 0

     go func() {
          for {
                select {
                      case ch <- fmt.Sprintf("%s,%d", str, i):
                            i++
                            //wg.Add(1)  //WG iS MEANT TO SYNC GR BEFORE ENTER GR
                      case <- quit:   //BLOCK GR if TRUE EXIT
                          return
                }
          }
     }()

     return ch
}

var wg sync.WaitGroup

func main() {
    quit := make(chan bool)

    //send
    ch :=  Sending("joe", quit)

    //receive
    wg.Add(20
    )
    go func() {
         for {
              select {
                  case s := <- ch:
                  fmt.Println(s)
                  wg.Done()
              }
         }
    }()


    wg.Wait()     //wait until sig
    quit <- true  //sig
   
     
    fmt.Println("End Program..!")
}