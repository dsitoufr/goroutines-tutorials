package main 

/////////////////////////////////////////
// How to user channels with a timeout 
//////////////////////////////////////////

import (
    "fmt"
    "time"
)

func watch(str string) <- chan string {
  
  ch := make(chan string) 

  go func() {
    for i := 0; i < 10; i++ {
        ch <- fmt.Sprintf("%s,%d", str, i)
        time.Sleep(time.Millisecond)
    }
  }()

  return  ch
}

func main() {

   ch := watch("joe")

    go func() {
        for {
            select {
                case s := <- ch: 
                     fmt.Println(s)
                //IF CANNOT READITERATION WITHIN TWO MILLI BREAKEXIT
                case <- time.After(2 * time.Millisecond) :
                     fmt.Println("This iteration is two slow I stop and exit")
                     return
            } 
        }
    }()
   
   time.Sleep(time.Second)
   fmt.Println("End Program.")
}