package shapes

import (
	"errors"
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

type Rectangle struct {
	Height, Width float64
}

func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

type Triangle struct {
	SideA, SideB, SideC float64
}

func (t Triangle) Area() float64 {
	p := (t.SideA + t.SideB + t.SideC) / 2
	return math.Sqrt(p * (p - t.SideA) * (p - t.SideB) * (p - t.SideC))
}

func ValidateTriangle(a, b, c float64) bool {
	return a <= (b+c) && b <= (a+c) && c <= (a+b)
}

func calculateArea(s any) (float64, error) {
	if shape, ok := s.(Shape); ok {
		return shape.Area(), nil
	}
	return 0.0, errors.New("переданный объект не реализует интерфейс Shape")
}

func roundToDecimal(n float64, r int) float64 {
	return math.Round(n*math.Pow10(r)) / math.Pow10(r)
}

func CalculateAreaCircle(radius float64) (float64, error) {
	circle := Circle{Radius: radius}
	if circleArea, err := calculateArea(circle); err == nil {
		fmt.Printf("Круг: радиус %.2f\nПлощадь: %.2f\n", radius, circleArea)
		fmt.Println()
		return roundToDecimal(circleArea, 2), nil
	} else {
		return 0.0, err
	}
}

func CalculateAreaRectangle(width, height float64) (float64, error) {
	rectangle := Rectangle{Width: width, Height: height}
	if rectangleArea, err := calculateArea(rectangle); err == nil {
		fmt.Printf("Прямоугольник: ширина %.2f, высота %.2f\nПлощадь: %.2f\n", width, height, rectangleArea)
		fmt.Println()
		return roundToDecimal(rectangleArea, 2), nil
	} else {
		return 0.0, err
	}
}

func CalculateAreaTriangle(a, b, c float64) (float64, error) {
	triangle := Triangle{SideA: a, SideB: b, SideC: c}
	if ValidateTriangle(a, b, c) {
		var s string
		if a == b && b == c {
			s = "равносторонний"
		}
		if a == b || b == c || a == c {
			s = "равнобедренный"
		}
		if triangleArea, err := calculateArea(triangle); err == nil {
			fmt.Printf("Треугольник: длины сторон %.2f, %.2f, %.2f, %v \n", a, b, c, s)
			fmt.Printf("Площадь: %.2f\n", triangleArea)
			fmt.Println()
			return roundToDecimal(triangleArea, 2), nil
		} else {
			return 0.0, err
		}
	} else {
		fmt.Printf("Треугольник: длины сторон %.2f, %.2f, %.2f \n\n"+
			"Для того чтобы треугольник существовал, сумма длин любых \n"+
			"двух его сторон должна быть больше длины третьей стороны.\n"+
			"Если это условие не выполняется, то треугольник невозможен.\n", a, b, c)
		return 0.0, fmt.Errorf("ошибка в передаваемых значениях")
	}
}
