package main 

import (
    "fmt"
    "time"
    "math/rand"
)

//RSS item
type Item struct { Title, Channel, GUID string}

//Fetch item
type Fetcher interface {
    Fetch() (items []Item, next time.Time, err error)
}

type Subscription interface {
    Updates() <- chan Item
    Close() error
}

func Subscribe(fetcher Fetcher) Subscription {
   s := &sub{
       fetcher: fetcher,
       updates: make(chan Item),  //for updates
       closing: make(chan chan error),  //for close
   }

   go s.loop()
   return s
}


type sub struct {
    fetcher Fetcher  
    updates chan  Item
    closing chan chan error
}

func(s *sub) Updates()  <-chan Item {
   return s.updates
}

func(s *sub) Close() error {

    errc := make(chan error)
    s.closing <- errc
    return <- errc
}

func(s *sub) loopFetchOnly() {
    var pending  []Item
    var next time.Time
    var err error

    for {
          var fetchDelay time.Duration
          if now := time.Now(); next.After(now) {
              fetchDelay = next.Sub(now)
          }

          startFetch := time.After(fetchDelay)

          select {
               case <-startFetch:
                    var  fetched []Item
                    fetched, next, err = s.fetcher.Fetch()
                    if err != nil {
                        next = time.Now().Add(10 * time.Second)
                        break
                    }
                    pending = append(pending, fetched...)
          }
    }
}


func (s *sub) loopSendOnly() {
  var pending []Item
  for {
        var first Item
        var updates chan Item
        if len(pending) > 0 {
             first = pending[0]
             updates = s.updates
        }

        select {
             case updates <- first:
                  pending = pending[1:]
        }
  }
} 


func (s *sub) mergedLoop() {
    var pending [] Item
    var next time.Time
    var err error

    for {
           var fetchDelay  time.Duration
           if now := time.Now(); next.After(now) {
                 fetchDelay = next.Sub(now)
           }
           startFetch := time.After(fetchDelay)
           var first Item
           var updates chan Item
           if len(pending) > 0 {
                first = pending[0]
                updates = s.updates
           }

           select {
               case   errc := <- s.closing:
                      errc <- err 
                      close(s.updates)
                      return
               case <- startFetch:
                     var fetched []Item
                     fetched, next, err = s.fetcher.Fetch()
                     if err != nil {
                          next = time.Now().Add(10 * time.Second)
                          break
                     }
                     pending = append(pending, fetched...)
                case updates <- first:
                       pending = pending[1:]
           }

    }
}


func(s *sub) dedupeLoop() {

    const maxPending = 10
    var pending []Item
    var next time.Time
    var err error
    var seen = make(map[string]bool)

    for {
          var fetchDelay time.Duration
          if now := time.Now(); next.After(now) {
               fetchDelay = next.Sub(now)
          }

          var startFetch <-chan time.Time
          if len(pending) < maxPending {
               startFetch = time.After(fetchDelay)
          }

          var first Item
          var updates chan Item

          if len(pending) > 0 {
               first = pending[0]
               updates = s.updates
          }

          select {
               case errc := <- s.closing:
                    errc <- err
                    close(s.updates)
                    return
               case  <- startFetch:
                     var fetched []Item
                     fetched, next, err = s.fetcher.Fetch()
                     if err != nil {
                         next = time.Now().Add(10 * time.Second)
                         break
                     }

                     for _, item := range fetched {
                          if !seen[item.GUID] {
                                pending = append(pending, item)
                                seen[item.GUID] = true
                          }
                     }
                case updates <- first:
                     pending = pending[1:]
          }

    }
}


func (s *sub) loop() {

    const maxPending = 10
    type fetchResult struct {
         fetched []Item
         next time.Time
         err error
    }

    var fetchDone chan fetchResult
    var pending []Item
    var next time.Time
    var err error
    var seen = make(map[string]bool)

    for {
         var fetchDelay time.Duration
         if now := time.Now(); next.After(now) {
             fetchDelay = next.Sub(now)
         }

         //STARTFETCH OMIT
         var startFetch <-chan time.Time
         if fetchDone == nil && len(pending) < maxPending {
               startFetch = time.After(fetchDelay)
         }

         // STOPFETCHIF OMIT
         var first Item
         var updates chan Item
         if len(pending)  > 0 {
             first = pending[0]
             updates = s.updates  //enable send case
         }

         //STARTFETCHASYNC OMIT
         select {
               case <-startFetch:
                    fetchDone = make(chan fetchResult,1)
                    go func() {
                         fetched, ext, err := s.fetcher.Fetch()
                         fetchDone <- fetchResult{fetched, next, err}
                    }()
                case result := <- fetchDone:
                    fetchDone = nil
                    fetched := result.fetched
                    next, err = result.next, result.err
                    if err != nil {
                        next = time.Now().Add(10 * time.Second)
                        break
                    }

                    for _, item := range fetched {
                        if id := item.GUID; !seen[id] {
                             pending = append(pending, item)
                             seen[id] = true
                        }
                    }
                case errc := <- s.closing:
                     errc <- err
                     close(s.updates)
                     return
                case  updates <- first:
                      pending = pending[1:] 
         }
    }
}

//STARTNAIVEMERGE OMIT

type naiveMerge struct {
  subs []Subscription
  updates chan Item
}

func NaiveMerge(subs ...Subscription) Subscription {
	m := &naiveMerge{
		subs:    subs,
		updates: make(chan Item),
	}
	for _, sub := range subs {
		go func(s Subscription) {
			for it := range s.Updates() {
				m.updates <- it // HL
			}
		}(sub)
	}
	return m
}

//STOPNAIVEMERGE  OMIT

//STARTNAIVEMERGECLOSE OMIT

func(m *naivMerge) Updates() <-chan Item {
      return m.updates
} 

type merge struct {
    subs []Subscription
    updates chan Item
    quit chan struct{}
    errs chan error
}

// STARTMERGESIG OMIT
func Merge(subs ...Subscription) Subscription {
    //STOPMERGESIG OMIT
    m := &merge{
        subs: subs,
        updates: make(chan Item),
        quit: make(chan struct{}),
        errs: make(chan error),
    }

    //STARTMERGE OMIT
    for _,sub := range subs {
          go func(s Subscription) {
              for {
                   var it Item
                   select {
                        case  it = <- s.updates():
                        case <- m.quit:
                          m.errs <- s.Close()
                          return
                   }

                   select {
                        case m.updates <- it:
                        case <- m.quit:
                                m.errs <- s.Close()
                                return
                   }
              }
          }(sub)
    }

    return m
}

func (* merge) Upadate() <- chan Item{
    return m.updates
}

//STARTMERGECLOSE OMIT
func NaiveDedupe(  in <-chan Item) {
    out := make(chan Item)
    go func() {
         seen := make(map[string]bool)
         for it := range in {
               if !seen[it.GUID] {
                   out <- it
                   seen[it.GUID] = true
               }
         }
         close(out)
    }()

    return out
}

type deduper struct {
  s Subscription
  updates chan Item
  closing chan chan error
}

func Dedupe(s Subscription) Subscription {
    d := &deduper{
         s:  s,
         updates: make(chan Item),
         closing: make(chan chan error),
    }
    go d.loop()
    return d
}

func(d *deduper) loop() {
    in := d.s.Updates()
    var pending Item
    var out chan Itemseen = make(map[String]bool)

    for {
         select {
             case it <- in:
                if !seen[it.GUID] {
                    pending = it
                    in = nil
                    out = d.updates
                    seen[it.GUID] = true
                }
            case out <- pending:
                 in = d.s.Updates()
                 out = nil
            case errc := <- d.closing:
                 err := d.s.close()
                 errc <- err
                 close(d.updates)
                 return
         }
    }

}

func(d *deduper) Close() error {
    errc := make(chan error)
    d.closing <- errc
    return ( <- errc)
}

func(d *deduper) Updates() <-chan Item {
    return d.updates
}

func Fetch(domain string) Fetcher {
      return fakeFetch(domain)
}

func fakeFetch(domain string) Fetcher {
    return &fakeFetcher{channel: domain}
}

type fakeFetcher struct {
     channel string
     items []Item
}

var FakeDuplicates bool

func(f *fakeFetcher) Fetch() (items []Item, next time.Time, err error) {

    now := time.Now()
    next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
    item := Item{
          Channel: f.channel,
          Title:  fmt.Sprintf("Item %d", len(f.items)),
    }

    item.GUID = item.Channel + "/" + item.Title
    f.items = append(f.items, item)
    if KakeDuplicates {
         items = f.items
    } else {
         items = []Item{item}
    }

    return
}

func init() {
     rand.Seed(time.Now().UnixNano())
}

//STARTMAIN OMIT

func main() {

    //Subscribes to some feeds and create merge stream
    merged := Merge(
        Subscribe(Fetch("blog.golang.org")),
        Subscribe(Fetch("googleblog.blogspot.com")),
        Subscribe(Fetch("googledevelopers.blogspot.com")))

        //STOPMERGECALL  OMIT
        //Close the subscriptions after some time
        time.AfterFunc("3*time.Second", func() {
             fmt.Println("closed:", merged.Close())
        })

        //Print the streams
        for it := range merged.Updates() {
              fmt.Println(it.Channel, it.Title)
        }

    panic("show me the stacks")
}