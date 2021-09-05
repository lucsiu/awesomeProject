package concurrency

import (
	"fmt"
	"time"
)

func Select() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		c2 <- "three"
	}()

	for i := 0; i < 3; i++ { // 因为需要获取三个chan的输出, 所以循环3次
		fmt.Printf("Select在这里阻塞住了, i:%v\n", i)
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

// 移除select外层for的方法1: 超时
func Timeouts() {
	// By default channels are unbuffered, meaning that they will only accept sends (chan <-)
	// if there is a corresponding receive (<- chan) ready to receive the sent value.
	//
	// Buffered channels accept a limited number of values without a corresponding receiver for those values.
	//
	// Note that the channel is buffered, so the send in the goroutine is nonblocking.
	c1 := make(chan string, 1)
	// external call
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
		fmt.Println("Note that the channel is buffered, so the send in the goroutine is nonblocking.")
	}()

	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()

	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}

// 移除select外层for的方法2: 无阻塞
func NonBlockingSelect() {
	//messages := make(chan string, 1) // bufferd, no-blocking
	messages := make(chan string) // no buffered, blocking
	signals := make(chan bool)

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	// Here msg cannot be sent to the messages channel,
	// because the channel has no buffer and there is no receiver.
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func TestChannelNoBuffered() {
	messages := make(chan string)

	go func() {
		messages <- "Greetings!"
		fmt.Println("'Greetings!' sent.") // never exec this, because the chann is not buffered.
	}()

	time.Sleep(2 * time.Second)

	result := <-messages
	fmt.Println(result)
}

func TestChannelBuffered() {
	messages := make(chan string, 1)

	go func() {
		messages <- "Greetings!"
		// Note that the channel is buffered, so the send in the goroutine is nonblocking.
		// will be exec.
		fmt.Println("'Greetings!' sent.")
	}()

	time.Sleep(2 * time.Second)

	result := <-messages
	fmt.Println(result)
}
