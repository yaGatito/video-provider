package shape

import "math"

type Shape interface {
	WithPerimeterShape
	WithAreaShape
}

type WithPerimeterShape interface {
	Perimeter() float64
}

type WithAreaShape interface {
	Area() float64
}

type Circle struct {
	radius float64
}

func NewCircle(radius float64) Circle {
	return Circle{radius: radius}
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius //  P = 2 * π * r
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius // A = π * r^2
}

type Square struct {
	side float64
}

func (s Square) Perimeter() float64 {
	return 4 * s.side // P = 4 * a
}

func (s Square) Area() float64 {
	return s.side * s.side // A = a^2
}
