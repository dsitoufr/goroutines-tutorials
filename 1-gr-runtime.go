package main

import (
"fmt"
"time"
)

func main() {
  for i := 0; i <= 9; i++ {
     i := i
      go func() {
          fmt.Println(i)
      }()
  }
  time.Sleep(time.Second)
}
