package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var siganChannel = make(chan os.Signal, 1)

func main() {
	//ContextWithCancel()
	ContextWithTimeout()
	//ContextWithDeadline()
}

func ContextWithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println(ctx.Value("howie"))
				return
			default:
				func() {
					fmt.Println("hello")
					time.Sleep(time.Second * 2)
					fmt.Println("hello")
				}()
			}
		}
	}()
	fmt.Println("开始")
	time.AfterFunc(5*time.Second, func() {
		ctx = context.WithValue(ctx, "howie", "10秒后调用cancel()")
		cancel()
		fmt.Println("结束")
	})
	Exit()

}

func ContextWithTimeout() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println(time.Since(start))
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				return
			}
		}
	}()
	Exit()

}

func ContextWithDeadline() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	start := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println(time.Since(start))
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			return
		}
	}()
	Exit()
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	// Associate the cancellable context we just created to the request
	req = req.WithContext(ctx)
	//client := &http.Client{}
	//res, err := client.Do(req)

}

func Exit() {
	signal.Notify(siganChannel, os.Kill, os.Interrupt)
	<-siganChannel
}
