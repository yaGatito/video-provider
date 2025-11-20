package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/zhashkevych/scheduler"
)

func main() {
	//generics.Start()
	main2()
}

func main2() {
	b := [5]float64{1.1, 2.2, 3.3, 4.4, 5.5}
	b[0] = 0.00001
	a := b
	fmt.Printf("%s\n", strconv.FormatFloat(a[0], 'f', -1, 64))
	fmt.Printf("%s\n", strconv.FormatFloat(b[0], 'f', -1, 64))

}

func main1() {
	ctx := context.Background()

	worker := scheduler.NewScheduler()
	worker.Add(ctx, parseLatestPosts, time.Second*5)
	worker.Add(ctx, collectStatistics, time.Second*10)

	time.AfterFunc(time.Minute*1, worker.Stop)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
}

func parseLatestPosts(ctx context.Context) {
	time.Sleep(time.Second * 1)
	fmt.Printf("latest posts parsed successfuly at %s\n", time.Now().String())
}

func collectStatistics(ctx context.Context) {
	time.Sleep(time.Second * 5)
	fmt.Printf("stats updated at %s\n", time.Now().String())
}
