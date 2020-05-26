package main

import (
"fmt"
)

func main() {
     extvar := "Hello "
     f := func(s string) {
             fmt.Println(extvar + s)
          } 
    f("roland")
    extvar = "Goodbye "
    f("roland")
}
