package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	urlTest = "amqp://guest:guest@localhost:5672/"
)

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	g, c := errgroup.WithContext(ctx)

	g.Go(func() error {
		d := make(chan int, 2)
		go func() {
			for i := 0; i < 10; i++ {
				d <- i
				time.Sleep(time.Second * 2)
			}
		}()
		for {
			select {
			case <-c.Done():
				fmt.Println("ctx done 1")
				return nil
			case i := <-d:
				fmt.Println("i:", i)
			}
		}
	})
	g.Go(func() error {
		select {
		case <-c.Done():
			fmt.Println("ctx done 2")
		}
		return nil
	})

	<-s
	cancelFunc()
	_ = g.Wait()
	//parse, err := url.Parse(urlTest)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v", parse)
}
