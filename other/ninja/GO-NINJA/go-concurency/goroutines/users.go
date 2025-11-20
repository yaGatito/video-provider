package goroutines

import (
	"fmt"
	"math/rand"
	"time"
)

var actions = []string{
	"logged in",
	"logged out",
	"create record",
	"delete record",
	"update record",
}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func saveUserInfo(user User) error {
	time.Sleep(time.Millisecond * 10)
	fmt.Println("Saving user: ", user.id)
	//file, err := os.OpenFile(fmt.Sprintf("%v.txt", user.id), os.O_RDWR|os.O_CREATE, 0644)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//
	//_, err = file.WriteString(user.getActivityInfo())
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	return nil
}

func generateUsers(count int) []User {
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@gmail.com", i+1),
			logs:  generateLogs(1000),
		}
	}

	return users
}

func generateUsers_withChannel(users chan User, count int) {
	for i := 0; i < count; i++ {
		users <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@gmail.com", i+1),
			logs:  generateLogs(1000),
		}
		fmt.Println("user generated and sent")
	}

	close(users)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			timestamp: time.Now(),
			action:    actions[rand.Intn(len(actions)-1)],
		}
	}

	return logs
}
