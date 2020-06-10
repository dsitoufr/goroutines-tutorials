//
// Daisy chain with only one single step
// One entry and one response
// This will serve as template for 
// multiply steps daisy chain
//
//

package main

import (
    "fmt"
)

func daisyOne(in , out chan int) {
    out <-  1+ <- in
}

func main() {
  out := make(chan int)
  in := make(chan int)

  go daisyOne(in ,out)
  in <- 1

  fmt.Println(<- out)
}
