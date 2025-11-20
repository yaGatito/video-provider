package channels

import (
	"fmt"
	"time"
)

func Start() {
	iterating_over_channel_elements()
}

// how to care channels
func example_channels() {
	fmt.Println("start app")

	// It seems deferred functions are not performing on fatal error occurred (like "all goroutines are asleep" - deadlock)
	defer func() {
		fmt.Println("deferred")
	}()

	msg := make(chan string)

	go func(msgChan chan string) {
		message := "pshol nahui"
		time.Sleep(time.Second * 2)
		fmt.Println("sentMessage message", message)
		msgChan <- message
	}(msg)

	go func(msgChan chan string) {
		// causing to sleep goroutine in order to wait for some resource
		select {
		case message := <-msgChan:
			time.Sleep(time.Second * 1)
			fmt.Println("handled message", message)
		}
	}(msg)

	time.Sleep(time.Second * 5)
}

// iterating over channel and its shortcut
func iterating_over_channel_elements() {
	msgs := make(chan string, 3)

	// waiting untill or all goroutines will be asleep to fatal error
	// OR until some goroutine(s) will set data to this channel
	msgs <- "asd"
	msgs <- "asd"
	msgs <- "asd"

	// deadlock if there is no next line
	close(msgs)

	// full loop over channel elements declaration
	for {
		value, ok := <-msgs
		if !ok {
			fmt.Println("channel closed")
			break
		}

		fmt.Println(value)
	}

	// shortcut loop over channel elements
	for m := range msgs {
		fmt.Println(m)
	}
}
