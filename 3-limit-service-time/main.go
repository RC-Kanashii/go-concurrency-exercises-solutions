//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	var start time.Time

	// 会员没有时间限制
	if u.IsPremium {
		start = time.Now()
		process()
		atomic.AddInt64(&u.TimeUsed, int64(time.Since(start)))
		fmt.Printf("User %v has used %vs\n", u.ID, float64(u.TimeUsed)/float64(time.Second))
		return true
	}

	// 超出时长限制，直接返回false，避免执行process()
	if atomic.LoadInt64(&u.TimeUsed) > int64(time.Second*10) {
		return false
	}

	resCh := make(chan bool)    // 传递执行结果
	timeUpCh := make(chan bool) // 传递是否超时

	// 该goroutine用于处理process的相关操作
	go func() {
		// 用来传递process()执行完毕信息的上下文
		ctx, cancel := context.WithCancel(context.Background())
		start = time.Now()
		// 单开一个goroutine，执行process
		go func() {
			process()
			// 告知父goroutine执行完毕
			cancel()
		}()
		// 轮询，监测是否达到用户使用时长限制，或者process是否执行完毕
		for {
			select {
			case <-timeUpCh:
				resCh <- false
				break
			case <-ctx.Done():
				atomic.AddInt64(&u.TimeUsed, int64(time.Since(start)))
				fmt.Printf("User %v has used %vs\n", u.ID, float64(u.TimeUsed)/float64(time.Second))
				resCh <- true
				break
			}
		}
	}()

	// 另开一个goroutine，实时监测用户是否达到时长限制
	go func() {
		for {
			if atomic.LoadInt64(&u.TimeUsed) > int64(time.Second*10) {
				// 通知上面的goroutine
				timeUpCh <- true
				close(timeUpCh)
				//fmt.Printf("User %v has used %vs. Time up.\n", u.ID, float64(u.TimeUsed)/float64(time.Second))
				break
			}
		}
	}()

	// 轮询
	for {
		select {
		case res := <-resCh:
			return res
		case <-time.After(time.Second * 10):
			// 此处仅增加时长，超时处理逻辑交由上面的goroutine处理
			atomic.AddInt64(&u.TimeUsed, int64(time.Second*10))
		}
	}
}

func main() {
	RunMockServer()
}
