package main

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"errors"
	"log"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	r float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.r, 2)
}

type Rectangle struct {
	h, w float64
}

func (r Rectangle) Area() float64 {
	return r.h * r.w
}

type Triangle struct {
	a, b, c float64
}

func (t *Triangle) TriangleSides() (float64, float64, float64) { // Getter
	return t.a, t.b, t.c
}

func (t *Triangle) SetTriangleSides(a, b, c float64) { // Setter
	t.a, t.b, t.c = a, b, c
}

func (t *Triangle) IsTriangle() bool {
	isTriangle := false
	if t.a <= (t.b+t.c) && t.b <= (t.a+t.c) && t.c <= (t.a+t.b) {
		isTriangle = true
	}
	return isTriangle
}

func (t Triangle) Area() float64 {
	a := 0.0
	if t.IsTriangle() {
		fmt.Printf("\nХарактеристика: ")
		if t.a == t.b && t.b == t.c {
			fmt.Println("равносторонний")
		} else if t.a == t.b || t.b == t.c || t.a == t.c {
			fmt.Println("равнобедренный")
		} else {
			fmt.Println("неравносторонний")
		}
		p := (t.a + t.b + t.c)                               // периметр
		s := p / 2                                           // полупериметр
		a = math.Sqrt(p * (s - t.a) * (s - t.b) * (s - t.c)) // площадь
	} else {
		fmt.Println("Какой-то не такой \"треугольник\"")
	}
	return a
}

type Square struct {
	a float64
}

func (s Square) calculateSquareArea() float64 { // здесь можно поменять на Area() и ошибки не будет
	return math.Pow(s.a, 2)
}

func calculateArea(s any) (float64, error) {
	if s, ok := s.(Shape); ok {
		return s.Area(), nil
	} else {
		return 0, errors.New("переданный объект не реализует интерфейс Shape")
	}
}

func getTypeInfo(s any) string {
	return reflect.TypeOf(s).String()
}

func main() {
	radiusValue := 5.6
	circleShape := Circle{radiusValue}
	if circleArea, err := calculateArea(circleShape); err == nil {
		fmt.Printf("Круг: радиус %.2f\nПлощадь: %.2f\n", radiusValue, circleArea)
	}
	fmt.Println()
	
	a := 5.2
	b := 6.4
	rectangleShape := Rectangle{a, b}
	if rectangleArea, err := calculateArea(rectangleShape); err == nil {
		fmt.Printf("Прямоугольник: ширина %.2f, высота %.2f\nПлощадь: %.2f\n", a, b, rectangleArea)
	}
	fmt.Println()

	a, b = 5, 7
	c := 10.1 // 12.1
	triangleShape := Triangle{a, b, c}
	fmt.Printf("Треугольник: длины сторон %.1f, %.1f, %.1f ", a, b, c)
	if !triangleShape.IsTriangle() {
		fmt.Printf("\nДля того чтобы треугольник существовал, сумма длин любых \n" +
				   "двух его сторон должна быть больше длины третьей стороны.\n" +
				   "Если это условие не выполняется, то треугольник невозможен.\n")

		fmt.Println("Для закрепления материала, воспользуемся Setter'ом, установим, что \n" +
		            "Треугольник: длины сторон 4.1, 5.2, 6")
		triangleShape.SetTriangleSides(4.1, 5.2, 6)
	}
	if triangleArea, err := calculateArea(triangleShape); err == nil {
		fmt.Printf("Площадь: %.2f\n", triangleArea)
	}
	fmt.Println()

	d := 10.0
	squareShape := Square{d}
	if _, err := calculateArea(squareShape); err != nil { // форма проверки
		log.Fatal(getTypeInfo(squareShape), " - ошибка: ", err)
	} else {
		fmt.Println(calculateArea(squareShape))
	}

	os.Exit(0)
}
