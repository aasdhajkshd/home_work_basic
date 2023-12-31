package shapes

import (
	"math"
)

type Shape interface {
	Area() float64
}

// Circle calculation
type Circle struct {
	Radius float64 // radius circle value
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

// Square for test purpose, none exportable, on via setter/getter
type Square struct {
	side float64 // square side value
}

func (s Square) CalculateSquareArea() float64 { // Due to hw05 may demonstrate a Shape interface error
	return math.Pow(s.side, 2)
}

// Getting the side of the Square
func (s *Square) SquareSide() float64 {
	return s.side
}

// Setting the side of the Square
func (s *Square) SetSquareSide(side float64) {
	s.side = side
}

// Rectangle calculation
type Rectangle struct {
	Height, Width float64
}

func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Triangle calculation
type Triangle struct {
	SideA, SideB, SideC float64
}

func (t Triangle) Area() float64 {
	p := (t.SideA + t.SideB + t.SideC)                                         // perimeter
	s := p / 2                                                                 // semiperimeter
	return math.Sqrt(p * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))        // area
}

func ValidateTriangle(a, b, c float64) bool {
	return a <= (b + c) && b <= (a + c) && c <= (a + b) == true                // validate traiable by sides' values
}
