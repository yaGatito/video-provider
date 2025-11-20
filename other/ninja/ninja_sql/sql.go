package main

import "fmt"

// docker run -d --name ninja-db -e POSTGRES_PASSWORD=1111 -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres
// docker exec -it ninja-db bash
// psql -U postgres
// https://www.udemy.com/course/golang-ninja/learn/lecture/34874084#content

func main() {
	fmt.Println("asd")
}
