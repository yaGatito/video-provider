package generics

type Number interface {
	int64 | float64
}

type User struct {
	Name  string
	Email string
}

func Start() {
	a := []int64{1, 2, 3, 4, 5}
	b := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	s := []string{"1.1", "2.2", "3.3", "4.4", "5.5"}

	users := []User{
		{Name: "Alice", Email: "alice@mail.com"},
		{Name: "Bob", Email: "bob@mail.com"},
		{Name: "Charlie", Email: "charlie@mail.com"},
	}

	sumInt := sum(a)
	sumFloat := sum(b)
	println("Sum of integers:", sumInt) // Output: Sum of integers: 15
	println("Sum of floats:", sumFloat) // Output: Sum of floats:

	println(searchElement(a, 3))
	println(searchElement(b, 4.5))
	println(searchElement(b, 4.4))
	println(searchElement(s, "5.5"))

	println("alice", searchElement(users, User{Name: "Alice", Email: "alice@mail.com"}))
	println("alic2e", searchElement(users, User{Name: "Alice", Email: "alic2e@mail.com"}))
}

func sum[V Number](input []V) V {
	var result V
	for _, number := range input {
		result += number
	}
	return result
}

func searchElement[C comparable](elements []C, searchEl C) bool {
	for _, v := range elements {
		if v == searchEl {
			return true
		}
	}
	return false
}
