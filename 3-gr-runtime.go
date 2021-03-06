package main

//
// You can interrupt a GR and force another one to execute
//

import (
"fmt"
"runtime"
)

func main() {
  var i byte
  go func() {
      for i := 0; i <= 255; i++ {
      }
  }()

  fmt.Println("First step ", i)
  runtime.Gosched()    //force other GR to execute
  runtime.GC()         //Gargage collector
  fmt.Println("Done")
}
