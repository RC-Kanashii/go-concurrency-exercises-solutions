//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, ch chan<- *Tweet, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if errors.Is(err, ErrEOF) {
			close(ch)
			return err
		}

		ch <- tweet
	}
}

func consumer(ch <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	ch := make(chan *Tweet, 10)
	wg := sync.WaitGroup{}

	// Producer
	wg.Add(1)
	go func() {
		err := producer(stream, ch, &wg)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Consumer
	wg.Add(1)
	go consumer(ch, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}

/*
Original Results:
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
vampirewalk666  tweets about golang
Process took 3.5824285s

After Optimization:
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
End of File
vampirewalk666  tweets about golang
Process took 1.979719667s
*/
