package main


import (
"fmt"
"time"
)

func count( s string ) {
  for i:=0; i < 5; i++ {
     fmt.Println(s)
     time.Sleep(100 * time.Millisecond)
  }
}

func main() {

 go count("hello")
 count("world")
 
 fmt.Println("End Program..!")
}
