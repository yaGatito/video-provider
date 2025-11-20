package goroutines

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func (u User) getActivityInfo() string {
	out := fmt.Sprintf("ID: %d\nEmail: %s\nLogs:\n", u.id, u.email)
	for i, item := range u.logs {
		out += fmt.Sprintf("%d. [%s] at %s", i, item.action, item.timestamp.Format("2006-01-02 15:04:05"))
	}

	return out
}

func Start() {
	t := time.Now()
	work1000goroutines()
	var firstMethodDeltaTime = time.Now().Sub(t)

	workGOMAXPROCSgoroutines()
	var secondMethodDeltaTime = time.Now().Sub(t) - firstMethodDeltaTime
	fmt.Println("first method duration: ", firstMethodDeltaTime)
	fmt.Println("second method duration: ", secondMethodDeltaTime)
}

func Start2() {
	users := make(chan User)
	// chanel closer method is next line
	go generateUsers_withChannel(users, 3)

	//fmt.Println(<-users)
	//fmt.Println(<-users)

	//for {
	//	fmt.Println("waiting for users channel")
	//	user, ok := <-users
	//	if !ok {
	//		fmt.Println("users channel closed")
	//		break
	//	}
	//	fmt.Println("received user with id", user.id)
	//	go saveUserInfo(user)
	//}

	for user := range users {
		fmt.Println("received user with id", user.id)
		go saveUserInfo(user)
	}

	users <- User{
		id:    0,
		email: "asd",
		logs: []logItem{logItem{
			action:    "asda",
			timestamp: time.Time{},
		}},
	}

	//for {
	//	fmt.Println("waiting for users channel")
	//	user, ok := <-users
	//	if !ok {
	//		fmt.Println("users channel closed")
	//		break
	//	}
	//	fmt.Println("received user with id", user.id)
	//	go saveUserInfo(user)
	//}

	//time.Sleep(100 * time.Second)
}

func work1000goroutines() {
	wg := &sync.WaitGroup{}

	users := generateUsers(1024)
	for _, user := range users {
		wg.Go(func() {
			err := saveUserInfo(user)
			if err != nil {
				return
			}
		})
	}

	wg.Wait()
}

// GOMAXPROCS
func workGOMAXPROCSgoroutines() {
	wg := &sync.WaitGroup{}

	// 16
	var goMax = runtime.GOMAXPROCS(runtime.NumCPU())
	var count = 1024
	var portion = count / goMax

	for i := 0; i < goMax; i++ {
		wg.Go(func() {
			users := generateUsers(portion)
			for _, user := range users {
				err := saveUserInfo(user)
				if err != nil {
					return
				}
			}
		})
	}

	wg.Wait()
}
