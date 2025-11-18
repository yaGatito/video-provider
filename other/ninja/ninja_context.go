package main

import (
	"context"
	"fmt"
	"time"
)

func main1() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "id", 1)

	parse(ctx)
}

func parse(ctx context.Context) {
	dTimeNow := time.Now()
	fmt.Println("dTimeNow:", dTimeNow)
	//  working on exact time stamp
	ctx, _ = context.WithDeadline(ctx, dTimeNow.Add(time.Second*3))
	// working with duration (timer duration)
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	value := ctx.Value("id")
	fmt.Println(value)
	for {
		select {
		case t := <-time.After(time.Second * 4):
			fmt.Println("parsing completed: ", t.String())
			return
		case <-ctx.Done():
			dTimeNow = time.Now()
			fmt.Println("dTimeDone:", dTimeNow)
			fmt.Println("deadline exceeded")
			return
		}
	}
}
