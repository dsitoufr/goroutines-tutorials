package main

import(
    "fmt"
    "time"
)

type Ball struct {
    hits int
}

func Player(name string, ch chan *Ball) {

   for {
        ball := <- ch
        ball.hits++
        fmt.Println(name, ball.hits)
        time.Sleep(time.Millisecond)
        ch <- ball
   }

}

func main() {
   ch := make( chan *Ball)
   // ch <-  new(Ball)  //Dead lock  the first send/receive to/from chan must be from goroutine
     

   go Player("ping", ch)
   go Player("pong", ch)

   ch <-  new(Ball)   //start on  
                      //new return an adress we need a pointer

   time.Sleep(10 * time.Millisecond)     //extend delay to have more exchange

   <- ch
   fmt.Println("End program.") 
   panic("show me the stack")

}